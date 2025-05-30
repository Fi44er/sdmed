package usecase

import (
	"context"
	"regexp"
	"time"

	"github.com/Fi44er/sdmedik/backend/internal/config"
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/dto"
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/entity"
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmedik/backend/internal/module/notification/service"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/Fi44er/sdmedik/backend/pkg/utils"
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

type AuthUsecase struct {
	logger         *logger.Logger
	cache          ICache
	config         *config.Config
	userUsecase    IUserUsecase
	notifyerServce *service.NotificationService
}

func NewAuthUsecase(
	logger *logger.Logger,
	cache ICache,
	config *config.Config,
	userUsecase IUserUsecase,
	notificationService *service.NotificationService,
) *AuthUsecase {
	return &AuthUsecase{
		logger:         logger,
		config:         config,
		cache:          cache,
		userUsecase:    userUsecase,
		notifyerServce: notificationService,
	}
}

const (
	CodeRedisPrefix = "verification_codes_"
	UserRedisPrefix = "temp_user_"
)

func (s *AuthUsecase) createAndStoreToken(ctx context.Context, userID string, expiresIn time.Duration, privateKey string, userAgent string) (string, error) {
	tokenDetails, err := utils.CreateToken(userID, expiresIn, privateKey)
	if err != nil {
		return "", constant.ErrUnprocessableEntity
	}
	key := userAgent + ":" + tokenDetails.TokenUUID
	err = s.cache.Set(ctx, key, userID, time.Until(time.Unix(*tokenDetails.ExpiresIn, 0)))
	return *tokenDetails.Token, err
}

func (s *AuthUsecase) SignIn(ctx context.Context, user *entity.User) (*entity.Tokens, error) {
	user, err := s.userUsecase.GetByEmail(ctx, user.Email)
	if err != nil || !utils.ComparePassword(user.Password, user.Password) {
		return nil, constant.ErrInvalidEmailOrPassword
	}

	accessToken, err := s.createAndStoreToken(ctx, user.ID, s.config.AccessTokenExpiresIn, s.config.AccessTokenPrivateKey, user.UserAgent)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.createAndStoreToken(ctx, user.ID, s.config.RefreshTokenExpiresIn, s.config.RefreshTokenPrivateKey, user.UserAgent)
	if err != nil {
		return nil, err
	}

	return &entity.Tokens{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthUsecase) VerifyCode(ctx context.Context, data *dto.VerifyCodeDTO) error {
	hashEmail, err := utils.HashString(data.Email)
	if err != nil {
		return constant.ErrInternalServerError
	}

	var code string
	if err := s.cache.Get(ctx, CodeRedisPrefix+hashEmail, &code); err != nil {
		return constant.ErrInternalServerError
	}

	if err := s.cache.Del(ctx, CodeRedisPrefix+hashEmail); err != nil {
		return err
	}

	var tempUser dto.SignUpDTO
	if err := s.cache.Get(ctx, UserRedisPrefix+hashEmail, &tempUser); err != nil {
		return err
	}

	return s.cache.Del(ctx, UserRedisPrefix+hashEmail)
}

func (s *AuthUsecase) SignUp(ctx context.Context, entity *entity.User) error {
	entity.PhoneNumber = regexp.MustCompile("[^0-9]").ReplaceAllString(entity.PhoneNumber, "")
	if len(entity.PhoneNumber) != 11 {
		return constant.ErrInvalidPhoneNumber
	}

	user, err := s.userUsecase.GetByEmail(ctx, entity.Email)
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

	if err := s.cache.Set(ctx, UserRedisPrefix+hashEmail, entity, 10*time.Minute); err != nil {
		return err
	}

	return s.SendCode(ctx, entity.Email)
}

func (s *AuthUsecase) SendCode(ctx context.Context, email string) error {
	code := utils.GenerateCode(6)
	hashEmail, err := utils.HashString(email)
	if err != nil {
		return err
	}

	var tempUser dto.SignUpDTO
	if err := s.cache.Get(ctx, UserRedisPrefix+hashEmail, &tempUser); err != nil {
		return constant.ErrUnprocessableEntity
	}

	if err := s.cache.Set(ctx, CodeRedisPrefix+hashEmail, code, s.config.VerifyCodeExpiredIn); err != nil {
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

	s.notifyerServce.Send(msg, "smtp")

	s.logger.Info(code)
	return nil
}

func (s *AuthUsecase) RefreshAccessToken(ctx context.Context, data *dto.RefreshTokenDTO) (string, error) {
	if data.RefreshToken == "" {
		return "", constant.ErrInvalidToken
	}

	tokenClaims, err := utils.ValidateToken(data.RefreshToken, s.config.RefreshTokenPublicKey)
	if err != nil {
		return "", constant.ErrForbidden
	}

	var userID string
	if err := s.cache.Get(ctx, tokenClaims.TokenUUID, &userID); err != nil {
		return "", constant.ErrCouldNotRefreshToken
	}

	user, err := s.userUsecase.GetByID(ctx, userID)
	if err != nil || user == nil {
		return "", err
	}

	return s.createAndStoreToken(ctx, user.ID, s.config.AccessTokenExpiresIn, s.config.AccessTokenPrivateKey, data.UserAgent)
}

func (s *AuthUsecase) SignOut(ctx context.Context, data *dto.LogoutDTO) error {
	if data.RefreshToken == "" {
		return constant.ErrInvalidToken
	}

	tokenClaims, err := utils.ValidateToken(data.RefreshToken, s.config.RefreshTokenPublicKey)
	if err != nil {
		return constant.ErrAnauthorized
	}

	if err := s.cache.Del(ctx, tokenClaims.TokenUUID); err != nil {
		return err
	}

	return nil
}
