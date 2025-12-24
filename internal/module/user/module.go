package user_module

import (
	user_handler "github.com/Fi44er/sdmed/internal/module/user/delivery/http/user"
	role_repository "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/role"
	user_repository "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/user"
	role_usecase "github.com/Fi44er/sdmed/internal/module/user/usecase/role"
	user_usecase "github.com/Fi44er/sdmed/internal/module/user/usecase/user"
	"github.com/gofiber/fiber/v2"

	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserModule struct {
	userRepository *user_repository.UserRepository
	userUsecase    *user_usecase.UserUsecase
	userHandler    *user_handler.UserHandler

	roleRepository role_repository.IRoleRepository
	roleUsecase    role_usecase.IRoleUsecase

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

	m.roleRepository = role_repository.NewRoleRepository(m.logger, m.db)
	m.roleUsecase = role_usecase.NewRoleUsecase(m.logger, m.roleRepository)
}

func (m *UserModule) InitDelivery(router fiber.Router) {
	m.userHandler.RegisterRoutes(router)
}

func (m *UserModule) GetUserUsecase() *user_usecase.UserUsecase {
	return m.userUsecase
}

func (m *UserModule) GetRoleUsecase() role_usecase.IRoleUsecase {
	return m.roleUsecase
}
