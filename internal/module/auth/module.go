package auth_module

import (
	"context"
	"log"

	"github.com/Fi44er/sdmed/internal/config"
	auth_handler "github.com/Fi44er/sdmed/internal/module/auth/delivery/http"
	auth_adapters "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/adapters"
	repository "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/session"
	user_session_repository "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/user_session"
	accessmanager_service "github.com/Fi44er/sdmed/internal/module/auth/usecase/access_manager"
	auth_usecase "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/shadow_user"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	role_usecase "github.com/Fi44er/sdmed/internal/module/user/usecase/role"
	user_usecase "github.com/Fi44er/sdmed/internal/module/user/usecase/user"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/redis"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthModule struct {
	authAdapters          *auth_adapters.UserUsecaseAdapter
	authUsecase           *auth_usecase.AuthUsecase
	authHandler           *auth_handler.AuthHandler
	userUsecase           *user_usecase.UserUsecase
	roleUsecase           role_usecase.IRoleUsecase
	notificationServce    *service.NotificationService
	sessionRepository     *repository.SessionRepository
	userSessionRepository user_session_repository.IUserSessionRepository
	shadowUserUsecase     shadow_user.IShadowUserService

	accessManager *accessmanager_service.Manager

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
	roleUsecase role_usecase.IRoleUsecase,
	notificationService *service.NotificationService,
) *AuthModule {
	return &AuthModule{
		logger:             logger,
		validator:          validator,
		db:                 db,
		redisManager:       redisManager,
		config:             config,
		userUsecase:        userUsecase,
		roleUsecase:        roleUsecase,
		notificationServce: notificationService,
	}
}

func (m *AuthModule) Init() {
	m.authAdapters = auth_adapters.NewUserUsecaseAdapter(m.userUsecase, m.roleUsecase)
	m.sessionRepository = repository.NewSessionRepository(m.logger)
	m.userSessionRepository = user_session_repository.NewUserSessionRepository(m.logger, m.db)
	m.shadowUserUsecase = shadow_user.NewShadowUserService(m.logger, m.authAdapters, m.config)
	m.authUsecase = auth_usecase.NewAuthUsecase(
		m.logger,
		m.redisManager,
		m.config,
		m.authAdapters,
		m.notificationServce,
		m.sessionRepository,
		m.userSessionRepository,
		m.shadowUserUsecase,
	)

	m.authHandler = auth_handler.NewAuthHandler(m.authUsecase, m.logger, m.validator, m.config)
	m.accessManager, _ = accessmanager_service.NewManager(m.db, m.authAdapters)

	if err := m.accessManager.SyncRolePermissions(context.Background()); err != nil {
		log.Fatal("Failed to sync policies:", err)
	}
}

func (m *AuthModule) InitDelivery(router fiber.Router) {
	m.authHandler.RegisterRoutes(router)
}

func (m *AuthModule) GetAccessManager() *accessmanager_service.Manager {
	return m.accessManager
}
