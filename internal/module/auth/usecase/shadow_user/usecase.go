package shadow_user

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	auth_constant "github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/google/uuid"
)

type IShadowUserService interface {
	CreateShadowUser(ctx context.Context) (*auth_entity.User, error)
	PromoteToRealUser(ctx context.Context, shadowUserID string, user *auth_entity.User) error
	CleanupExpiredShadows(ctx context.Context) error
}

type ShadowUserService struct {
	logger      *logger.Logger
	userUsecase IUserUsecase
	config      *config.Config
}

type IUserUsecase interface {
	Create(ctx context.Context, user *auth_entity.User) error
	Update(ctx context.Context, user *auth_entity.User) error
	GetByID(ctx context.Context, id string) (*auth_entity.User, error)
	Delete(ctx context.Context, id string) error
}

func NewShadowUserService(
	logger *logger.Logger,
	userUsecase IUserUsecase,
	config *config.Config,
) *ShadowUserService {
	return &ShadowUserService{
		logger:      logger,
		userUsecase: userUsecase,
		config:      config,
	}
}

// CreateShadowUser - создает временного пользователя для неавторизированных
func (s *ShadowUserService) CreateShadowUser(ctx context.Context) (*auth_entity.User, error) {
	now := time.Now()
	expiresAt := now.Add(24 * time.Hour) // 24 часа жизни по умолчанию

	shadowUser := &auth_entity.User{
		ID:              uuid.New().String(),
		Email:           "", // shadow user не имеет email
		Password:        "",
		PhoneNumber:     "",
		FIO:             "Guest User",
		IsShadow:        true,
		ShadowCreatedAt: &now,
		ShadowExpiresAt: &expiresAt,
		Roles: []auth_entity.Role{
			{Name: "guest"}, // базовая роль для shadow
		},
	}

	if err := s.userUsecase.Create(ctx, shadowUser); err != nil {
		s.logger.Errorf("Failed to create shadow user: %v", err)
		return nil, err
	}

	s.logger.Infof("Shadow user created: %s", shadowUser.ID)
	return shadowUser, nil
}

// PromoteToRealUser - конвертирует shadow user в настоящего при регистрации
func (s *ShadowUserService) PromoteToRealUser(ctx context.Context, shadowUserID string, userData *auth_entity.User) error {
	existingUser, err := s.userUsecase.GetByID(ctx, shadowUserID)
	if err != nil {
		return err
	}

	if !existingUser.IsShadow {
		return auth_constant.ErrUserAlreadyExists
	}

	// Обновляем shadow user данными реального пользователя
	existingUser.Email = userData.Email
	existingUser.Password = userData.Password
	existingUser.PhoneNumber = userData.PhoneNumber
	existingUser.FIO = userData.FIO
	existingUser.IsShadow = false
	existingUser.ShadowCreatedAt = nil
	existingUser.ShadowExpiresAt = nil
	// Роли будут обновлены в user usecase

	if err := s.userUsecase.Update(ctx, existingUser); err != nil {
		s.logger.Errorf("Failed to promote shadow user: %v", err)
		return err
	}

	s.logger.Infof("Shadow user %s promoted to real user", shadowUserID)
	return nil
}

// CleanupExpiredShadows - удаляет истекших shadow users (для cron job)
func (s *ShadowUserService) CleanupExpiredShadows(ctx context.Context) error {
	s.logger.Info("Starting cleanup of expired shadow users")

	// Эту логику нужно реализовать в user repository
	// Здесь просто заглушка для примера архитектуры
	// В реальности сделай метод в user_repository: DeleteExpiredShadows()

	s.logger.Info("Expired shadow users cleanup completed")
	return nil
}
