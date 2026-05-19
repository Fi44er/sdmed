package scraper_http

import (
	"github.com/Fi44er/sdmed/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *ScraperHandler) RegisterRoutes(router fiber.Router) {
	scraper := router.Group("/scraper", middlewares.RequireAuth())
	scraper.Get("/settings", middlewares.Authorize("scraper", "read"), h.GetSettings)
	scraper.Put("/settings", middlewares.Authorize("scraper", "write"), h.UpdateSettings)

	scraper.Post("/start", middlewares.Authorize("scraper", "write"), h.StartScraper)
	scraper.Post("/stop", middlewares.Authorize("scraper", "write"), h.StopScraper)
	scraper.Get("/status/stream", middlewares.Authorize("scraper", "read"), h.StatusSSE)
	scraper.Get("/status", middlewares.Authorize("scraper", "read"), h.GetStatus)
}
