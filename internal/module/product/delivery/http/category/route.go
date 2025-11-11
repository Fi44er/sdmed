package category_http

import "github.com/gofiber/fiber/v2"

func (h *CategoryHandler) RegisterRoutes(router fiber.Router) {
	categories := router.Group("/categories")
	categories.Post("/", h.Create)
}
