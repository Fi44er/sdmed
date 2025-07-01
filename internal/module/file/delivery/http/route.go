package http

import "github.com/gofiber/fiber/v2"

func (h *FileHandler) RegisterRoutes(router fiber.Router) {
	file := router.Group("/files")
	file.Post("/upload", h.Upload)
	file.Get("/:name", h.Get)
}
