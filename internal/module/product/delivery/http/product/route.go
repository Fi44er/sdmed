package product_http

import "github.com/gofiber/fiber/v2"

func (h *ProductHandler) RegisterRoutes(router fiber.Router) {
	products := router.Group("/products")
	products.Post("/", h.Create)
	products.Get("/:slug", h.GetBySlug)
}
