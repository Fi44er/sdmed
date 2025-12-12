package category_http

import "github.com/gofiber/fiber/v2"

func (h *CategoryHandler) RegisterRoutes(router fiber.Router) {
	categories := router.Group("/categories")
	categories.Post("/", h.Create)
	categories.Get("/:id", h.GetByID)
	categories.Get("/", h.GetAll)
	categories.Delete("/:id", h.Delete)
	categories.Put("/:id", h.Update)
	categories.Get("/by-slug/:slug", h.GetBySlug)

	// categories.Get("/filters/:category_id", h.GetFiltersByCategoryID)
}
