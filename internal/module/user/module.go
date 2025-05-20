package user

import (
	user_handler "github.com/Fi44er/sdmedik/backend/internal/module/user/delivery/http/user"
	user_repository "github.com/Fi44er/sdmedik/backend/internal/module/user/infrastructure/repository/user"
	user_usecase "github.com/Fi44er/sdmedik/backend/internal/module/user/usecase/user"
	"github.com/gofiber/fiber/v2"

	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserModule struct {
	userRepository *user_repository.UserRepository
	userUsecase    *user_usecase.UserUsecase
	userHandler    *user_handler.UserHandler

	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB
}

func NewUserModule(
	logger *logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
) *UserModule {
	return &UserModule{
		logger:    logger,
		validator: validator,
		db:        db,
	}
}

func (m *UserModule) Init() {
	m.userRepository = user_repository.NewUserRepository(m.logger, m.db)
	m.userUsecase = user_usecase.NewUserUsecase(m.userRepository, m.logger)
	m.userHandler = user_handler.NewUserHandler(m.userUsecase, m.logger, m.validator)
}

func (m *UserModule) InitDelivery(router fiber.Router) {
	m.userHandler.RegisterRoutes(router)
}
