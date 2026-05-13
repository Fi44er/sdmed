package scraper_http

import "github.com/gofiber/fiber/v2"

func (h *ScraperHandler) RegisterRoutes(router fiber.Router) {
	scraper := router.Group("/scraper")
	scraper.Get("/settings", h.GetSettings)
	scraper.Put("/settings", h.UpdateSettings)

	scraper.Post("/start", h.StartScraper)
	scraper.Post("/stop", h.StopScraper)
	scraper.Get("/status/stream", h.StatusSSE)
	scraper.Get("/status", h.GetStatus)
}
