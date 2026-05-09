package scraper_module

import (
	scraper_http "github.com/Fi44er/sdmed/internal/module/scraper/delivery/http"
	scraper_settings_repository "github.com/Fi44er/sdmed/internal/module/scraper/infrastructure/repository/scraper_settings"
	ratelimiter_service "github.com/Fi44er/sdmed/internal/module/scraper/service/rate_limiter"
	scraper_service "github.com/Fi44er/sdmed/internal/module/scraper/service/scraper"
	scrapersettings_service "github.com/Fi44er/sdmed/internal/module/scraper/service/scraper_settings"
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
	scraperHandler         *scraper_http.ScraperHandler

	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB
}

func NewScraperModule(logger *logger.Logger, validator *validator.Validate, db *gorm.DB) *ScraperModule {
	return &ScraperModule{
		logger:    logger,
		validator: validator,
		db:        db,
	}
}

func (m *ScraperModule) Init() {
	m.scraperSettingRepo = scraper_settings_repository.NewScraperSettingsRepository(m.db, m.logger)
	m.rateLimiterService = ratelimiter_service.NewRateLimiterService(m.logger, m.scraperSettingRepo)
	m.scraperService = scraper_service.NewScraperService(m.logger, m.rateLimiterService)
	m.scraperSettingsService = scrapersettings_service.NewScraperSettingsService(m.logger, m.scraperSettingRepo)
	m.scraperHandler = scraper_http.NewScraperHandler(m.scraperSettingsService, m.scraperService, m.logger, m.validator)
}

func (m *ScraperModule) InitDelivery(router fiber.Router) {
	m.scraperHandler.RegisterRoutes(router)
}
