package order_module

import (
	cart_usecase "github.com/Fi44er/sdmed/internal/module/cart/usecase/cart"
	notification_service "github.com/Fi44er/sdmed/internal/module/notification/service"
	order_http "github.com/Fi44er/sdmed/internal/module/order/delivery/http"
	order_adapters "github.com/Fi44er/sdmed/internal/module/order/infrastructure/adapters"
	order_payment "github.com/Fi44er/sdmed/internal/module/order/infrastructure/payment"
	order_repository "github.com/Fi44er/sdmed/internal/module/order/infrastructure/repository/order"
	order_usecase "github.com/Fi44er/sdmed/internal/module/order/usecase/order"
	product_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/product"
	tru_usecase "github.com/Fi44er/sdmed/internal/module/tru/usecase"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderModule struct {
	repo    order_repository.IOrderRepository
	usecase *order_usecase.OrderUsecase
	handler *order_http.OrderHandler

	db                  *gorm.DB
	logger              *logger.Logger
	validator           *validator.Validate
	config              *config.Config
	cartUsecase         *cart_usecase.CartUsecase
	productUsecase      product_usecase.IProductUsecase
	truUsecase          tru_usecase.ITRUCodeUseCase
	notificationService *notification_service.NotificationService
}

func NewOrderModule(
	db *gorm.DB,
	logger *logger.Logger,
	validator *validator.Validate,
	config *config.Config,
	cartUsecase *cart_usecase.CartUsecase,
	productUsecase product_usecase.IProductUsecase,
	truUsecase tru_usecase.ITRUCodeUseCase,
	notificationService *notification_service.NotificationService,
) *OrderModule {
	return &OrderModule{
		db:                  db,
		logger:              logger,
		validator:           validator,
		config:              config,
		cartUsecase:         cartUsecase,
		productUsecase:      productUsecase,
		truUsecase:          truUsecase,
		notificationService: notificationService,
	}
}

func (m *OrderModule) Init() {
	m.repo = order_repository.NewOrderRepository(m.db, m.logger)

	paymentProvider := order_payment.NewPayKeeperProvider(m.config, m.logger)

	cartAdapter := order_adapters.NewCartAdapter(m.cartUsecase)
	productAdapter := order_adapters.NewProductAdapter(m.productUsecase)
	truAdapter := order_adapters.NewTRUAdapter(m.truUsecase)
	notificationAdapter := order_adapters.NewNotificationAdapter(m.notificationService)
	chatAdapter := order_adapters.NewChatAdapter()

	m.usecase = order_usecase.NewOrderUsecase(
		m.repo,
		paymentProvider,
		cartAdapter,
		productAdapter,
		truAdapter,
		notificationAdapter,
		chatAdapter,
		m.logger,
		m.validator,
		m.config,
	)

	m.handler = order_http.NewOrderHandler(m.usecase, m.logger, m.validator)
}

func (m *OrderModule) InitDelivery(router fiber.Router, authMiddleware, adminMiddleware fiber.Handler) {
	m.handler.RegisterRoutes(router, authMiddleware, adminMiddleware)
}
