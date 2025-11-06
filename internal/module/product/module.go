package product_module

import (
	file_usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	product_adapters "github.com/Fi44er/sdmed/internal/module/product/infrastucture/adapters"
	category_repository "github.com/Fi44er/sdmed/internal/module/product/infrastucture/repository/category"
	category_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/category"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductModule struct {
	categoryRepository category_repository.ICategoryRepository
	categoryUsecase    category_usecase.ICategoryUsecase

	fileUsecaseAdapter product_adapters.IFileUsecaseAdapter
	fileUsecase        file_usecase.IFileUsecase

	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB
}

func NewProductModule(
	logger *logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
) *ProductModule {
	return &ProductModule{
		logger:    logger,
		validator: validator,
		db:        db,
	}
}

func (m *ProductModule) Init() {
	m.fileUsecaseAdapter = product_adapters.NewFileUsecaseAdapter(m.fileUsecase)
	m.categoryRepository = category_repository.NewCategoryRepository(m.logger, m.db)
	m.categoryUsecase = category_usecase.NewCategoryUsecase(m.logger, m.categoryRepository, m.fileUsecaseAdapter)

}

func (m *ProductModule) InitDelivery(router fiber.Route) {
}
