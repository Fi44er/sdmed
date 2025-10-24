package auth_module

import (
	"github.com/Fi44er/sdmed/internal/config"
	auth_handler "github.com/Fi44er/sdmed/internal/module/auth/delivery/http"
	"github.com/Fi44er/sdmed/internal/module/auth/infrastucture/adapters"
	repository "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/session"
	auth_usecase "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	user_usecase "github.com/Fi44er/sdmed/internal/module/user/usecase/user"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/redis"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthModule struct {
	authauth_adapters  *auth_adapters.UserUsecaseAdapter
	authUsecase        *auth_usecase.AuthUsecase
	authHandler        *auth_handler.AuthHandler
	userUsecase        *user_usecase.UserUsecase
	notificationServce *service.NotificationService
	sessionRepository  *repository.SessionRepository
	tokenService       *auth_adapters.TokenService

	logger       *logger.Logger
	validator    *validator.Validate
	db           *gorm.DB
	redisManager redis.IRedisManager
	config       *config.Config
}

func NewAuthModule(
	logger *logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
	redisManager redis.IRedisManager,
	config *config.Config,
	userUsecase *user_usecase.UserUsecase,
	notificationService *service.NotificationService,
) *AuthModule {
	return &AuthModule{
		logger:             logger,
		validator:          validator,
		db:                 db,
		redisManager:       redisManager,
		config:             config,
		userUsecase:        userUsecase,
		notificationServce: notificationService,
	}
}

func (m *AuthModule) Init() {
	m.authauth_adapters = auth_adapters.NewUserUsecaseAdapter(m.userUsecase)
	m.sessionRepository = repository.NewSessionRepository(m.logger)
	m.tokenService = auth_adapters.NewTokenService()
	m.authUsecase = auth_usecase.NewAuthUsecase(
		m.logger,
		m.redisManager,
		m.config,
		m.authauth_adapters,
		m.notificationServce,
		m.sessionRepository,
		m.tokenService,
	)
	m.authHandler = auth_handler.NewAuthHandler(m.authUsecase, m.logger, m.validator, m.config)
}

func (m *AuthModule) InitDelivery(router fiber.Router) {
	m.authHandler.RegisterRoutes(router)
}
