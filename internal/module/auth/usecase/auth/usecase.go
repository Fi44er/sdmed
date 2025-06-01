package usecase

import (
	"context"
	"regexp"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/auth/dto"
	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/utils"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type IUserUsecase interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key string) error
}

type INotificationService interface {
	Send(msg *service.Message, selectedNotifiers ...string)
}

type ISessionRepository interface {
	GetSessionInfo(ctx context.Context) (*entity.UserSesion, error)
	PutSessionInfo(ctx context.Context, sessionInfo *entity.UserSesion) error
	DeleteSessionInfo(ctx context.Context) error
}

type AuthUsecase struct {
	logger            *logger.Logger
	cache             ICache
	config            *config.Config
	userUsecase       IUserUsecase
	notifyerService   INotificationService
	sessionRepository ISessionRepository
}

func NewAuthUsecase(
	logger *logger.Logger,
	cache ICache,
	config *config.Config,
	userUsecase IUserUsecase,
	notificationService INotificationService,
	sessionRepository ISessionRepository,
) *AuthUsecase {
	return &AuthUsecase{
		logger:            logger,
		config:            config,
		cache:             cache,
		userUsecase:       userUsecase,
		notifyerService:   notificationService,
		sessionRepository: sessionRepository,
	}
}

const (
	CodeRedisPrefix = "verification_codes_"
	UserRedisPrefix = "temp_user_"
)

func (u *AuthUsecase) createToken(userID string, expiresIn time.Duration, privateKey string) (string, error) {
	tokenDetails, err := utils.CreateToken(userID, expiresIn, privateKey)
	if err != nil {
		return "", constant.ErrUnprocessableEntity
	}

	return *tokenDetails.Token, err
}

func (u *AuthUsecase) SignIn(ctx context.Context, user *entity.User) (*entity.Tokens, error) {
	existingUser, err := u.userUsecase.GetByEmail(ctx, user.Email)
	if err != nil || !utils.ComparePassword(existingUser.Password, user.Password) {
		return nil, constant.ErrInvalidEmailOrPassword
	}

	accessToken, err := u.createToken(user.ID, u.config.AccessTokenExpiresIn, u.config.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.createToken(user.ID, u.config.RefreshTokenExpiresIn, u.config.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	userSession := &entity.UserSesion{
		UserID:       existingUser.ID,
		RefreshToken: refreshToken,
	}

	if err := u.sessionRepository.PutSessionInfo(ctx, userSession); err != nil {
		return nil, err
	}

	return &entity.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (u *AuthUsecase) VerifyCode(ctx context.Context, verifyCode *entity.VerifyCode) error {
	hashEmail, err := utils.HashString(verifyCode.Email)
	if err != nil {
		return constant.ErrInternalServerError
	}

	var code string
	if err := u.cache.Get(ctx, CodeRedisPrefix+hashEmail, &code); err != nil {
		return constant.ErrInternalServerError
	}

	if err := u.cache.Del(ctx, CodeRedisPrefix+hashEmail); err != nil {
		return err
	}

	var tempUser dto.SignUpDTO
	if err := u.cache.Get(ctx, UserRedisPrefix+hashEmail, &tempUser); err != nil {
		return err
	}

	user := &entity.User{
		PhoneNumber: tempUser.PhoneNumber,
		Email:       tempUser.Email,
		Password:    tempUser.Password,
		FIO:         tempUser.FIO,
	}

	if err := u.userUsecase.Create(ctx, user); err != nil {
		return err
	}

	return u.cache.Del(ctx, UserRedisPrefix+hashEmail)
}

func (u *AuthUsecase) SignUp(ctx context.Context, entity *entity.User) error {
	entity.PhoneNumber = regexp.MustCompile("[^0-9]").ReplaceAllString(entity.PhoneNumber, "")
	if len(entity.PhoneNumber) != 11 {
		return constant.ErrInvalidPhoneNumber
	}

	user, err := u.userUsecase.GetByEmail(ctx, entity.Email)
	if err != nil {

	}

	if user != nil {
		return constant.ErrUserAlreadyExists
	}

	entity.Password = utils.GeneratePassword(entity.Password)

	hashEmail, err := utils.HashString(entity.Email)
	if err != nil {
		return err
	}

	if err := u.cache.Set(ctx, UserRedisPrefix+hashEmail, entity, 10*time.Minute); err != nil {
		return err
	}

	return u.SendCode(ctx, entity.Email)
}

func (u *AuthUsecase) SendCode(ctx context.Context, email string) error {
	code := utils.GenerateCode(6)
	hashEmail, err := utils.HashString(email)
	if err != nil {
		return err
	}

	var tempUser dto.SignUpDTO
	if err := u.cache.Get(ctx, UserRedisPrefix+hashEmail, &tempUser); err != nil {
		return constant.ErrUnprocessableEntity
	}

	if err := u.cache.Set(ctx, CodeRedisPrefix+hashEmail, code, u.config.VerifyCodeExpiredIn); err != nil {
		return err
	}

	date := time.Now().Format("2006")
	templateData := struct {
		VerifyCode string
		Date       string
	}{
		VerifyCode: code,
		Date:       date,
	}

	msg := &service.Message{
		Recipient:    email,
		Subject:      "Код подтверждения регистрации",
		Data:         templateData,
		TemplatePath: "./internal/module/auth/pkg/template/verify_code.html",
	}

	u.notifyerService.Send(msg, "smtp")

	u.logger.Info(code)
	return nil
}

func (u *AuthUsecase) RefreshAccessToken(ctx context.Context) (string, error) {

	sessionInfo, err := u.sessionRepository.GetSessionInfo(ctx)
	if err != nil {
		return "", err
	}

	_, err = utils.ValidateToken(sessionInfo.RefreshToken, u.config.RefreshTokenPublicKey)
	if err != nil {
		return "", constant.ErrForbidden
	}

	user, err := u.userUsecase.GetByID(ctx, sessionInfo.UserID)
	if err != nil || user == nil {
		return "", err
	}

	return u.createToken(user.ID, u.config.AccessTokenExpiresIn, u.config.AccessTokenPrivateKey)
}

func (u *AuthUsecase) SignOut(ctx context.Context) error {
	if err := u.sessionRepository.DeleteSessionInfo(ctx); err != nil {
		return err
	}

	return nil
}
