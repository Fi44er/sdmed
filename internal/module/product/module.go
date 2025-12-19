package product_module

import (
	"github.com/Fi44er/sdmed/internal/config"
	file_usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	category_http "github.com/Fi44er/sdmed/internal/module/product/delivery/http/category"
	product_http "github.com/Fi44er/sdmed/internal/module/product/delivery/http/product"
	product_adapters "github.com/Fi44er/sdmed/internal/module/product/infrastructure/adapters"
	category_repository "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/category"
	characteristic_repository "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/characteristic"
	char_value_repository "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/characteristic_value"
	product_repository "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/product"
	category_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/category"
	char_value_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/char_value"
	characteristic_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/characteristic"
	product_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/product"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductModule struct {
	categoryRepository category_repository.ICategoryRepository
	categoryUsecase    category_usecase.ICategoryUsecase
	categoryHandler    *category_http.CategoryHandler

	fileUsecaseAdapter product_adapters.IFileUsecaseAdapter
	fileUsecase        file_usecase.IFileUsecase

	characteristicRepository characteristic_repository.ICharacteristicRepository
	characteristicUsecase    characteristic_usecase.ICharacteristicUsecase

	productRepository product_repository.IProductRepository
	productUsecase    product_usecase.IProductUsecase
	productHandler    *product_http.ProductHandler

	charValueRepository char_value_repository.ICharValueRepository
	charValueUsecase    char_value_usecase.ICharValueUsecase

	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB
	uow       uow.Uow
	config    *config.Config
}

func NewProductModule(
	logger *logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
	uow uow.Uow,
	fileUsecase file_usecase.IFileUsecase,
	config *config.Config,
) *ProductModule {
	return &ProductModule{
		logger:      logger,
		validator:   validator,
		db:          db,
		uow:         uow,
		fileUsecase: fileUsecase,
		config:      config,
	}
}

func (m *ProductModule) Init() {
	m.uow.RegisterRepository("category", func(tx *gorm.DB) (any, error) {
		return category_repository.NewCategoryRepository(m.logger, tx), nil
	})

	m.uow.RegisterRepository("characteristic", func(tx *gorm.DB) (any, error) {
		return characteristic_repository.NewCharacteristicRepository(m.logger, tx), nil
	})

	m.uow.RegisterRepository("product", func(tx *gorm.DB) (any, error) {
		return product_repository.NewProductRepository(m.logger, tx), nil
	})

	m.uow.RegisterRepository("char_value", func(tx *gorm.DB) (any, error) {
		return char_value_repository.NewCharValueRepository(m.logger, tx), nil
	})

	m.characteristicRepository = characteristic_repository.NewCharacteristicRepository(m.logger, m.db)
	m.characteristicUsecase = characteristic_usecase.NewCharacteristicUsecase(m.characteristicRepository, m.uow, m.logger)

	m.fileUsecaseAdapter = product_adapters.NewFileUsecaseAdapter(m.fileUsecase)
	m.categoryRepository = category_repository.NewCategoryRepository(m.logger, m.db)
	m.categoryUsecase = category_usecase.NewCategoryUsecase(m.logger, m.categoryRepository, m.fileUsecaseAdapter, m.characteristicUsecase, m.uow)
	m.categoryHandler = category_http.NewCategoryHandler(m.categoryUsecase, m.logger, m.validator, m.config)

	m.charValueRepository = char_value_repository.NewCharValueRepository(m.logger, m.db)
	m.charValueUsecase = char_value_usecase.NewCharValueUsecase(m.logger, m.charValueRepository, m.uow, m.characteristicUsecase)

	m.productRepository = product_repository.NewProductRepository(m.logger, m.db)
	m.productUsecase = product_usecase.NewProductUsecase(m.productRepository, m.logger, m.uow, m.fileUsecaseAdapter, m.charValueUsecase)
	m.productHandler = product_http.NewProductHandler(m.productUsecase, m.validator, m.logger, m.config)
}

func (m *ProductModule) InitDelivery(router fiber.Router) {
	m.categoryHandler.RegisterRoutes(router)
	m.productHandler.RegisterRoutes(router)
}
