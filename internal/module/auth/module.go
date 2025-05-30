package module

import (
	"github.com/Fi44er/sdmedik/backend/internal/config"
	auth_handler "github.com/Fi44er/sdmedik/backend/internal/module/auth/delivery/http"
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/infrastucture/adapters"
	auth_usecase "github.com/Fi44er/sdmedik/backend/internal/module/auth/usecase/auth"
	"github.com/Fi44er/sdmedik/backend/internal/module/notification/service"
	user_usecase "github.com/Fi44er/sdmedik/backend/internal/module/user/usecase/user"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/Fi44er/sdmedik/backend/pkg/redis"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type AuthModule struct {
	authAdapters       *adapters.UserUsecaseAdapter
	authUsecase        *auth_usecase.AuthUsecase
	authHandler        *auth_handler.AuthHandler
	userUsecase        *user_usecase.UserUsecase
	notificationServce *service.NotificationService

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
	m.authAdapters = adapters.NewUserUsecaseAdapter(m.userUsecase)
	m.authUsecase = auth_usecase.NewAuthUsecase(
		m.logger,
		m.redisManager,
		m.config,
		m.authAdapters,
		m.notificationServce,
	)
	m.authHandler = auth_handler.NewAuthHandler(m.authUsecase, m.logger, m.validator)
}

func (m *AuthModule) InitDelivery(router fiber.Router) {
	log.Warn(m.authAdapters == nil)
	m.authHandler.RegisterRoutes(router)
}
