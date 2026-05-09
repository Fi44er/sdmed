package scraper_http

import (
	"context"

	scraper_dto "github.com/Fi44er/sdmed/internal/module/scraper/dto"
	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IScraperSettingsService interface {
	Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error
	Get(ctx context.Context) (*scraper_entity.ScraperSettings, error)
}

type IScraperService interface {
	Scraper(ctx context.Context, onProgress func(completed, total int, message string)) []scraper_entity.Items
}

type ScraperHandler struct {
	settingsService IScraperSettingsService
	scraperService  IScraperService
	logger          *logger.Logger
	validator       *validator.Validate

	converter *Converter
}

func NewScraperHandler(settingsService IScraperSettingsService, scraperService IScraperService, logger *logger.Logger, validator *validator.Validate) *ScraperHandler {
	return &ScraperHandler{
		settingsService: settingsService,
		scraperService:  scraperService,
		logger:          logger,
		validator:       validator,
		converter:       &Converter{},
	}
}

// GetSettings godoc
// @Summary Get scraper settings
// @Tags Scraper
// @Accept json
// @Produce json
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /scraper/settings [get]
func (h *ScraperHandler) GetSettings(ctx *fiber.Ctx) error {
	settings, err := h.settingsService.Get(ctx.Context())
	if err != nil {
		h.logger.Errorf("error while getting settings: %s", err)
		return err
	}

	return ctx.Status(200).JSON(settings)
}

// UpdateSettings godoc
// @Summary Update scraper settings
// @Tags Scraper
// @Accept json
// @Produce json
// @Param body body scraper_dto.UpdateSettingsDTO true "Update Settings"
// @Success 200 {object} response.Response "OK"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /scraper/settings [put]
func (h *ScraperHandler) UpdateSettings(ctx *fiber.Ctx) error {
	dto := new(scraper_dto.UpdateSettingsDTO)

	// Используем ваш вспомогательный метод (предположим, конвертер ToEntityUpdateSettings)
	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntity, h.logger)
	if err != nil {
		return err
	}

	if err := h.settingsService.Update(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error while updating settings: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "settings updated successfully",
	})
}

// RunScraper godoc
// @Summary Run scraping process
// @Description Starts the scraper and returns results. Note: progress is logged server-side.
// @Tags Scraper
// @Accept json
// @Produce json
// @Success 200 {array} scraper_entity.Items "OK"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /scraper/run [post]
func (h *ScraperHandler) RunScraper(ctx *fiber.Ctx) error {
	// Callback для логирования прогресса
	onProgress := func(completed, total int, message string) {
		h.logger.Debugf("Scraping progress: %d/%d - %s", completed, total, message)
	}

	// Запуск скрапера
	results := h.scraperService.Scraper(ctx.Context(), onProgress)

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   results,
	})
}
