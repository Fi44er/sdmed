package cart_module

import (
	cart_http "github.com/Fi44er/sdmed/internal/module/cart/delivery/http"
	cart_adapters "github.com/Fi44er/sdmed/internal/module/cart/infrastructure/adapters"
	cart_repository "github.com/Fi44er/sdmed/internal/module/cart/infrastructure/repository/cart"
	cart_usecase "github.com/Fi44er/sdmed/internal/module/cart/usecase/cart"
	product_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/product"
	tru_usecase "github.com/Fi44er/sdmed/internal/module/tru/usecase"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartModule struct {
	repo           cart_repository.ICartRepository
	usecase        *cart_usecase.CartUsecase
	handler        *cart_http.CartHandler
	productAdapter *cart_adapters.ProductUsecaseAdapter
	truAdapter     *cart_adapters.TRUCodeUsecaseAdapter
	authMigrator   *cart_adapters.AuthMigrator

	db             *gorm.DB
	logger         *logger.Logger
	validator      *validator.Validate
	uow            uow.Uow
	productUsecase product_usecase.IProductUsecase
	truUsecase     tru_usecase.ITRUCodeUseCase
}

func NewCartModule(
	db *gorm.DB,
	logger *logger.Logger,
	validator *validator.Validate,
	uow uow.Uow,
	productUsecase product_usecase.IProductUsecase,
	truUsecase tru_usecase.ITRUCodeUseCase,
) *CartModule {
	return &CartModule{
		db:             db,
		logger:         logger,
		validator:      validator,
		uow:            uow,
		productUsecase: productUsecase,
		truUsecase:     truUsecase,
	}
}

func (m *CartModule) Init() {
	// Register repository in UOW if needed
	m.uow.RegisterRepository("cart", func(tx *gorm.DB) (any, error) {
		return cart_repository.NewCartRepository(tx, m.logger), nil
	})

	m.repo = cart_repository.NewCartRepository(m.db, m.logger)
	m.productAdapter = cart_adapters.NewProductUsecaseAdapter(m.productUsecase)
	m.truAdapter = cart_adapters.NewTRUCodeUsecaseAdapter(m.truUsecase)

	// Promo adapter could be injected here if available
	m.usecase = cart_usecase.NewCartUsecase(m.repo, m.productAdapter, m.truAdapter, nil, m.logger, m.validator)
	m.authMigrator = cart_adapters.NewAuthMigrator(m.usecase)
	m.handler = cart_http.NewCartHandler(m.usecase, m.logger, m.validator)
}

func (m *CartModule) InitDelivery(router fiber.Router) {
	m.handler.RegisterRoutes(router)
}

func (m *CartModule) GetUsecase() *cart_usecase.CartUsecase {
	return m.usecase
}

func (m *CartModule) GetCartMigrator() *cart_adapters.AuthMigrator {
	return m.authMigrator
}
