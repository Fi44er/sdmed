package product_http

import (
	"github.com/Fi44er/sdmed/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *ProductHandler) RegisterRoutes(router fiber.Router) {
	products := router.Group("/products")
	products.Post("/", middlewares.RequireAuth(), middlewares.Authorize("products", "write"), h.Create)
	products.Get("/:slug", h.GetBySlug)
	products.Get("/", h.GetAll)
	products.Get("/filters/:category_id", h.GetFilters)
}
