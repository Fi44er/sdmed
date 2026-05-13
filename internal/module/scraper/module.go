package scraper_module

import (
	product_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/product"
	scraper_http "github.com/Fi44er/sdmed/internal/module/scraper/delivery/http"
	scraper_adapters "github.com/Fi44er/sdmed/internal/module/scraper/infrastructure/adapters"
	scraper_settings_repository "github.com/Fi44er/sdmed/internal/module/scraper/infrastructure/repository/scraper_settings"
	ratelimiter_service "github.com/Fi44er/sdmed/internal/module/scraper/service/rate_limiter"
	scraper_service "github.com/Fi44er/sdmed/internal/module/scraper/service/scraper"
	scrapermanager_service "github.com/Fi44er/sdmed/internal/module/scraper/service/scraper_manager"
	scrapersettings_service "github.com/Fi44er/sdmed/internal/module/scraper/service/scraper_settings"
	tru_usecase "github.com/Fi44er/sdmed/internal/module/tru/usecase"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ScraperModule struct {
	scraperSettingRepo     scraper_settings_repository.IScraperSettingsRepository
	rateLimiterService     ratelimiter_service.IRateLimiterService
	scraperService         *scraper_service.ScraperService
	scraperSettingsService scrapersettings_service.IScraperSettingsService
	scraperManagerService  *scrapermanager_service.ScraperManagerService
	scraperHandler         *scraper_http.ScraperHandler

	productUseCase product_usecase.IProductUsecase
	truCodeUseCase tru_usecase.ITRUCodeUseCase

	productUseCaseAdapter *scraper_adapters.ProductUsecaseAdapter
	truCodeUseCaseAdapter *scraper_adapters.TRUUseCaseAdapter

	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB
}

func NewScraperModule(
	logger *logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
	productUseCase product_usecase.IProductUsecase,
	truCodeUseCase tru_usecase.ITRUCodeUseCase,
) *ScraperModule {
	return &ScraperModule{
		logger:         logger,
		validator:      validator,
		db:             db,
		productUseCase: productUseCase,
		truCodeUseCase: truCodeUseCase,
	}
}

func (m *ScraperModule) Init() {
	m.productUseCaseAdapter = scraper_adapters.NewProductUsecaseAdapter(m.productUseCase)
	m.truCodeUseCaseAdapter = scraper_adapters.NewTRUUseCaseAdapter(m.truCodeUseCase)

	m.scraperSettingRepo = scraper_settings_repository.NewScraperSettingsRepository(m.db, m.logger)
	m.rateLimiterService = ratelimiter_service.NewRateLimiterService(m.logger, m.scraperSettingRepo)
	m.scraperService = scraper_service.NewScraperService(m.logger, m.rateLimiterService)
	m.scraperSettingsService = scrapersettings_service.NewScraperSettingsService(m.logger, m.scraperSettingRepo)
	m.scraperManagerService = scrapermanager_service.NewScraperManagerService(m.logger, m.scraperService, m.productUseCaseAdapter, m.truCodeUseCaseAdapter)
	m.scraperHandler = scraper_http.NewScraperHandler(m.scraperSettingsService, m.scraperManagerService, m.logger, m.validator)
}

func (m *ScraperModule) InitDelivery(router fiber.Router) {
	m.scraperHandler.RegisterRoutes(router)
}
