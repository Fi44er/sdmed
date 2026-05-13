package scraper_http

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	scraper_dto "github.com/Fi44er/sdmed/internal/module/scraper/dto"
	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	scrapermanager_service "github.com/Fi44er/sdmed/internal/module/scraper/service/scraper_manager"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IScraperSettingsService interface {
	Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error
	Get(ctx context.Context) (*scraper_entity.ScraperSettings, error)
}

type IScraperManagerService interface {
	Start(ctx context.Context, params scraper_entity.ScrapeParams) error
	Stop()
	GetStatus() scrapermanager_service.ScraperStatusDTO
	Subscribe() chan scrapermanager_service.ScraperStatusDTO
	Unsubscribe(ch chan scrapermanager_service.ScraperStatusDTO)
}

type ScraperHandler struct {
	settingsService IScraperSettingsService
	managerService  IScraperManagerService
	logger          *logger.Logger
	validator       *validator.Validate

	converter *Converter
}

func NewScraperHandler(settingsService IScraperSettingsService, managerService IScraperManagerService, logger *logger.Logger, validator *validator.Validate) *ScraperHandler {
	return &ScraperHandler{
		settingsService: settingsService,
		managerService:  managerService,
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

// StartScraper godoc
// @Summary      Запустить процесс парсинга
// @Description  Запускает фоновый процесс сбора данных КТСР и синхронизации с БД.
// @Tags         Scraper
// @Accept       json
// @Produce      json
// @Param        request  body      StartScraperRequest  true  "Фильтры парсинга (оставьте пустыми для полного парсинга)"
// @Success      202      {object}  map[string]string    "Процесс запущен"
// @Failure      400      {object}  map[string]string    "Парсер уже запущен"
// @Router       /scraper/start [post]
func (h *ScraperHandler) StartScraper(c *fiber.Ctx) error {
	req := new(StartScraperRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Маппим DTO во внутренние параметры скрапера
	params := scraper_entity.ScrapeParams{
		Categories: req.Categories,
		Articles:   req.Articles,
		Regions:    req.Regions,
	}

	// Запускаем процесс.
	// Мы передаем context.Background(), чтобы процесс жил после завершения HTTP-ответа.
	if err := h.managerService.Start(context.Background(), params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Scraping process started in background",
	})
}

// StopScraper godoc
// @Summary      Остановить парсинг
// @Description  Прерывает выполнение текущей задачи парсинга.
// @Tags         Scraper
// @Produce      json
// @Success      200      {object}  map[string]string
// @Router       /scraper/stop [post]
func (h *ScraperHandler) StopScraper(c *fiber.Ctx) error {
	h.managerService.Stop()
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Stop signal sent",
	})
}

// StatusSSE godoc
// @Summary      Стрим статуса парсинга (SSE)
// @Description  Удерживает соединение и отправляет обновления прогресса в реальном времени
// @Tags         Scraper
// @Router       /scraper/status/stream [get]
func (h *ScraperHandler) StatusSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	// Создаем индивидуальную подписку
	statusChan := h.managerService.Subscribe()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// Обязательно отписываемся при выходе из функции (когда клиент отключится)
		defer h.managerService.Unsubscribe(statusChan)

		h.logger.Info("SSE: Client connected")

		for {
			// Ждем статус или используем таймер для "пинка" (keep-alive)
			select {
			case status, ok := <-statusChan:
				if !ok {
					return
				} // Канал закрыт сервером

				data, _ := json.Marshal(status)
				fmt.Fprintf(w, "data: %s\n\n", data)

				// ЕСЛИ Flush вернул ошибку — значит клиент ушел
				if err := w.Flush(); err != nil {
					h.logger.Info("SSE: Client disconnected (flush error)")
					return
				}
			case <-time.After(5 * time.Second):
				// Шлем пустой комментарий для поддержания связи
				fmt.Fprintf(w, ": keep-alive\n\n")
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	})

	return nil
}

// GetStatus godoc
// @Summary      Получить статус парсинга
// @Description  Возвращает информацию о текущем прогрессе, фазе и запущен ли процесс.
// @Tags         Scraper
// @Produce      json
// @Success      200      {object}  ScraperStatusResponse
// @Router       /scraper/status [get]
func (h *ScraperHandler) GetStatus(c *fiber.Ctx) error {
	status := h.managerService.GetStatus()
	return c.JSON(status)
}

type RunScraperRequest struct {
	Categories []string `json:"categories" example:"06,07"`
	Articles   []string `json:"articles" example:"06-01-01,07-02-01"`
	Regions    []string `json:"regions" example:"77,78"` // ISO коды регионов
}

type StartScraperRequest struct {
	Categories []string `json:"categories" example:"06,07"`
	Articles   []string `json:"articles" example:"06-01-01"`
	Regions    []string `json:"regions" example:"77,78"`
}

// ScraperStatusResponse — то, что мы отдаем в Swagger
type ScraperStatusResponse struct {
	IsRunning      bool    `json:"is_running"`
	CompletedTasks int     `json:"completed_tasks"`
	TotalTasks     int     `json:"total_tasks"`
	Percent        float64 `json:"percent"`
	Message        string  `json:"message"`
}
