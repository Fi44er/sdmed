package usecase

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/utils"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/contracts"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type AuthUsecase struct {
	logger            *logger.Logger
	cache             contracts.ICache
	config            *config.Config
	userUsecase       contracts.IUserUsecase
	notifyerService   contracts.INotificationService
	sessionRepository contracts.ISessionRepository
	tokenService      contracts.ITokenService
}

func NewAuthUsecase(
	logger *logger.Logger,
	cache contracts.ICache,
	config *config.Config,
	userUsecase contracts.IUserUsecase,
	notificationService contracts.INotificationService,
	sessionRepository contracts.ISessionRepository,
	tokenService contracts.ITokenService,
) *AuthUsecase {
	return &AuthUsecase{
		logger:            logger,
		config:            config,
		cache:             cache,
		userUsecase:       userUsecase,
		notifyerService:   notificationService,
		sessionRepository: sessionRepository,
		tokenService:      tokenService,
	}
}

const (
	CodeRedisPrefix = "verification_codes_"
	UserRedisPrefix = "temp_user_"
)

func (u *AuthUsecase) createToken(userID string, expiresIn time.Duration, privateKey string) (string, error) {
	tokenDetails, err := u.tokenService.CreateToken(userID, expiresIn, privateKey)
	if err != nil {
		return "", constant.ErrUnprocessableEntity
	}

	return *tokenDetails.Token, err
}

func (u *AuthUsecase) SignIn(ctx context.Context, user *entity.User) (*entity.Tokens, error) {
	existingUser, err := u.userUsecase.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if !u.userUsecase.ComparePassword(existingUser, user.Password) {
		return nil, constant.ErrInvalidEmailOrPassword
	}

	accessToken, err := u.createToken(existingUser.ID, u.config.AccessTokenExpiresIn, u.config.AccessTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.createToken(existingUser.ID, u.config.RefreshTokenExpiresIn, u.config.RefreshTokenPrivateKey)
	if err != nil {
		return nil, err
	}

	userSession := &entity.UserSession{
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

	var user entity.User
	if err := u.cache.Get(ctx, UserRedisPrefix+hashEmail, &user); err != nil {
		return err
	}

	if err := u.userUsecase.Create(ctx, &user); err != nil {
		return err
	}

	return u.cache.Del(ctx, UserRedisPrefix+hashEmail)
}

func (u *AuthUsecase) SignUp(ctx context.Context, user *entity.User) error {
	// user.PhoneNumber = regexp.MustCompile("[^0-9]").ReplaceAllString(user.PhoneNumber, "")
	if len(user.PhoneNumber) != 11 {
		return constant.ErrInvalidPhoneNumber
	}

	existUser, err := u.userUsecase.GetByEmail(ctx, user.Email)
	if err != nil {
		if err != constant.ErrUserNotFound {
			return err
		}
	}

	if existUser != nil {
		return constant.ErrUserAlreadyExists
	}

	user.Password = utils.GeneratePassword(user.Password)

	hashEmail, err := utils.HashString(user.Email)
	if err != nil {
		return err
	}

	if err := u.cache.Set(ctx, UserRedisPrefix+hashEmail, user, 10*time.Minute); err != nil {
		return err
	}

	sendCode := &entity.SendCode{
		Email: user.Email,
	}

	return u.SendCode(ctx, sendCode)
}

func (u *AuthUsecase) SendCode(ctx context.Context, sendCode *entity.SendCode) error {
	code := utils.GenerateCode(6)
	hashEmail, err := utils.HashString(sendCode.Email)
	if err != nil {
		return err
	}

	var tempUser entity.User
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
		Recipient:    sendCode.Email,
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

	_, err = u.tokenService.ValidateToken(sessionInfo.RefreshToken, u.config.RefreshTokenPublicKey)
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
