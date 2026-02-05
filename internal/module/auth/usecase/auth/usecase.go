package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	auth_constant "github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	auth_utils "github.com/Fi44er/sdmed/internal/module/auth/pkg/utils"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/contracts"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/google/uuid"
)

type IAuthUsecase interface {
	CreateShadowSession(ctx context.Context) (*auth_entity.User, error)
	SignIn(ctx context.Context, user *auth_entity.User) error
	VerifyCode(ctx context.Context, verifyCode *auth_entity.Code) error
	SignUp(ctx context.Context, user *auth_entity.User) error
	SendCode(ctx context.Context, sendCode *auth_entity.Code) error
	SignOut(ctx context.Context) error
	SignOutAll(ctx context.Context) error
	GetUserDevices(ctx context.Context) ([]*auth_entity.DeviceInfo, error)
	RevokeDevice(ctx context.Context, deviceID string) error
	ForgotPassword(ctx context.Context, code *auth_entity.Code) error
	ValidateResetPassword(ctx context.Context, token string) (string, error)
	ResetPassword(ctx context.Context, token string, user *auth_entity.User) error
}

type AuthUsecase struct {
	logger                *logger.Logger
	cache                 contracts.ICache
	config                *config.Config
	userUsecase           contracts.IUserUsecase
	notifyerService       contracts.INotificationService
	sessionRepository     contracts.ISessionRepository
	userSessionRepository contracts.IUserSessionRepository
	shadowUserService     contracts.IShadowUserService
}

func NewAuthUsecase(
	logger *logger.Logger,
	cache contracts.ICache,
	config *config.Config,
	userUsecase contracts.IUserUsecase,
	notificationService contracts.INotificationService,
	sessionRepository contracts.ISessionRepository,
	userSessionRepository contracts.IUserSessionRepository,
	shadowUserService contracts.IShadowUserService,
) *AuthUsecase {
	return &AuthUsecase{
		logger:                logger,
		config:                config,
		cache:                 cache,
		userUsecase:           userUsecase,
		notifyerService:       notificationService,
		sessionRepository:     sessionRepository,
		userSessionRepository: userSessionRepository,
		shadowUserService:     shadowUserService,
	}
}

const (
	CodeRedisPrefix           = "verification_codes_"
	UserRedisPrefix           = "temp_user_"
	ForgotPasswordRedisPrefix = "forgot_password_"
)

func (u *AuthUsecase) CreateShadowSession(ctx context.Context) (*auth_entity.User, error) {
	// Создаем shadow user
	shadowUser, err := u.shadowUserService.CreateShadowUser(ctx)
	if err != nil {
		return nil, err
	}

	deviceID := uuid.New().String()

	userSession := &auth_entity.ActiveSession{
		UserID:    shadowUser.ID,
		DeviceID:  deviceID,
		UserRoles: []string{"guest"},
		IsShadow:  true,
	}

	if err := u.sessionRepository.PutSessionInfo(ctx, userSession); err != nil {
		return nil, err
	}

	// Сохраняем в БД для управления устройствами
	dbSession := &auth_entity.UserSession{
		ID:         deviceID,
		UserID:     shadowUser.ID,
		DeviceName: "Guest Device",
		ExpiresAt:  time.Now().Add(u.config.RefreshTokenExpiresIn),
		IsRevoked:  false,
	}

	if err := u.userSessionRepository.Create(ctx, dbSession); err != nil {
		return nil, err
	}

	return shadowUser, nil
}

func (u *AuthUsecase) SignIn(ctx context.Context, user *auth_entity.User) error {
	existingUser, err := u.userUsecase.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if !u.userUsecase.ComparePassword(existingUser, user.Password) {
		return auth_constant.ErrInvalidEmailOrPassword
	}

	deviceID := uuid.New().String()

	strinRoles := make([]string, 0)
	for _, role := range existingUser.Roles {
		strinRoles = append(strinRoles, role.Name)
	}

	userSession := &auth_entity.ActiveSession{
		UserID:    existingUser.ID,
		DeviceID:  deviceID,
		UserRoles: strinRoles,
		IsShadow:  false,
	}

	if err := u.sessionRepository.PutSessionInfo(ctx, userSession); err != nil {
		return err
	}

	dbSession := &auth_entity.UserSession{
		ID:         deviceID,
		UserID:     existingUser.ID,
		DeviceName: "Unknown Device", // TODO: Парсить из User-Agent
		ExpiresAt:  time.Now().Add(u.config.RefreshTokenExpiresIn),
		IsRevoked:  false,
	}

	if err := u.userSessionRepository.Create(ctx, dbSession); err != nil {
		return err
	}

	return nil
}

func (u *AuthUsecase) VerifyCode(ctx context.Context, verifyCode *auth_entity.Code) error {
	hashEmail, err := auth_utils.HashString(verifyCode.Email)
	if err != nil {
		return auth_constant.ErrInternalServerError
	}

	var code string
	if err := u.cache.Get(ctx, CodeRedisPrefix+hashEmail, &code); err != nil {
		return auth_constant.ErrInternalServerError
	}

	if err := u.cache.Del(ctx, CodeRedisPrefix+hashEmail); err != nil {
		return err
	}

	var user auth_entity.User
	if err := u.cache.Get(ctx, UserRedisPrefix+hashEmail, &user); err != nil {
		return err
	}

	sessionInfo, err := u.sessionRepository.GetSessionInfo(ctx)
	u.logger.Debugf("session info in VerifyCode: %+v \n %v", sessionInfo, err)
	if err == nil && sessionInfo.IsShadow {
		// Конвертируем shadow user в реального
		if err := u.shadowUserService.PromoteToRealUser(ctx, sessionInfo.UserID, &user); err != nil {
			return err
		}
	} else {
		// Создаем нового пользователя
		if err := u.userUsecase.Create(ctx, &user); err != nil {
			return err
		}
	}

	return u.cache.Del(ctx, UserRedisPrefix+hashEmail)
}

func (u *AuthUsecase) SignUp(ctx context.Context, user *auth_entity.User) error {
	if len(user.PhoneNumber) != 11 {
		return auth_constant.ErrInvalidPhoneNumber
	}

	existUser, err := u.userUsecase.GetByEmail(ctx, user.Email)
	if err != nil {
		if err != auth_constant.ErrUserNotFound {
			return err
		}
	}

	if existUser != nil {
		return auth_constant.ErrUserAlreadyExists
	}

	user.Password = auth_utils.GeneratePassword(user.Password)

	hashEmail, err := auth_utils.HashString(user.Email)
	if err != nil {
		return err
	}

	if err := u.cache.Set(ctx, UserRedisPrefix+hashEmail, user, 10*time.Minute); err != nil {
		return err
	}

	sendCode := &auth_entity.Code{
		Email: user.Email,
		Type:  auth_entity.CodeTypeVerify,
	}

	return u.SendCode(ctx, sendCode)
}

func (u *AuthUsecase) SendCode(ctx context.Context, sendCode *auth_entity.Code) error {
	code := auth_utils.GenerateCode(6)
	hashEmail, err := auth_utils.HashString(sendCode.Email)
	if err != nil {
		return err
	}

	var tempUser auth_entity.User
	if err := u.cache.Get(ctx, UserRedisPrefix+hashEmail, &tempUser); err != nil {
		return auth_constant.ErrUnprocessableEntity
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

func (u *AuthUsecase) SignOut(ctx context.Context) error {
	sessionInfo, err := u.sessionRepository.GetSessionInfo(ctx)
	if err != nil {
		return err
	}

	if err := u.sessionRepository.DeleteSessionInfo(ctx); err != nil {
		return err
	}

	if err := u.userSessionRepository.Delete(ctx, sessionInfo.DeviceID); err != nil {
		u.logger.Warnf("Failed to delete session from DB: %v", err)
	}

	return nil
}

func (u *AuthUsecase) SignOutAll(ctx context.Context) error {
	sessionInfo, err := u.sessionRepository.GetSessionInfo(ctx)
	if err != nil {
		return err
	}

	return u.userSessionRepository.RevokeAllExcept(ctx, sessionInfo.UserID, sessionInfo.DeviceID)
}

func (u *AuthUsecase) GetUserDevices(ctx context.Context) ([]*auth_entity.DeviceInfo, error) {
	sessionInfo, err := u.sessionRepository.GetSessionInfo(ctx)
	if err != nil {
		return nil, err
	}

	sessions, err := u.userSessionRepository.GetByUserID(ctx, sessionInfo.UserID)
	if err != nil {
		return nil, err
	}

	devices := make([]*auth_entity.DeviceInfo, len(sessions))
	for i, session := range sessions {
		devices[i] = &auth_entity.DeviceInfo{
			DeviceID:   session.ID,
			DeviceName: session.DeviceName,
			UserAgent:  session.UserAgent,
			LastIP:     session.LastIP,
			IsCurrent:  session.ID == sessionInfo.DeviceID,
			CreatedAt:  session.CreatedAt,
			LastUsedAt: session.UpdatedAt,
		}
	}

	return devices, nil
}

func (u *AuthUsecase) RevokeDevice(ctx context.Context, deviceID string) error {
	sessionInfo, err := u.sessionRepository.GetSessionInfo(ctx)
	if err != nil {
		return err
	}

	// Проверяем, что устройство принадлежит этому пользователю
	session, err := u.userSessionRepository.Get(ctx, deviceID)
	if err != nil {
		return err
	}

	if session.UserID != sessionInfo.UserID {
		return auth_constant.ErrForbidden
	}

	return u.userSessionRepository.RevokeSession(ctx, deviceID)
}

func (u *AuthUsecase) ForgotPassword(ctx context.Context, code *auth_entity.Code) error {
	existUser, err := u.userUsecase.GetByEmail(ctx, code.Email)
	if err != nil {
		return err
	}

	token, err := auth_utils.GenerateSecretToken(32)
	if err != nil {
		return err
	}

	if err := u.cache.Set(ctx, ForgotPasswordRedisPrefix+token, existUser.ID, u.config.ResetPassTokenExpiredIn); err != nil {
		return err
	}

	resetLink := fmt.Sprintf("%s?token=%s", u.config.ResetPassURL, token)
	date := time.Now().Format("2006")
	templateData := struct {
		ResetLink string
		Date      string
	}{
		ResetLink: resetLink,
		Date:      date,
	}
	msg := &service.Message{
		Recipient:    code.Email,
		Subject:      "Сброс пароля",
		Data:         templateData,
		TemplatePath: "./internal/module/auth/pkg/template/reset_pass.html",
	}

	u.notifyerService.Send(msg, "smtp")

	return nil
}

func (u *AuthUsecase) ValidateResetPassword(ctx context.Context, token string) (string, error) {
	var userID string
	if err := u.cache.Get(ctx, ForgotPasswordRedisPrefix+token, &userID); err != nil {
		return "", err
	}
	if userID == "" {
		return "", auth_constant.ErrUnprocessableEntity
	}
	return userID, nil
}

func (u *AuthUsecase) ResetPassword(ctx context.Context, token string, user *auth_entity.User) error {
	userID, err := u.ValidateResetPassword(ctx, token)
	if err != nil {
		return err
	}
	user.Password = auth_utils.GeneratePassword(user.Password)
	user.ID = userID

	if err := u.userUsecase.Update(ctx, user); err != nil {
		return err
	}

	if err := u.cache.Del(ctx, ForgotPasswordRedisPrefix+token); err != nil {
		return err
	}

	if err := u.userSessionRepository.RevokeAll(ctx, userID); err != nil {
		u.logger.Warnf("Failed to revoke all sessions after password reset: %v", err)
	}

	return nil
}
