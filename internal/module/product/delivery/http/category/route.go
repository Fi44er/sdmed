package category_http

import (
	"github.com/Fi44er/sdmed/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *CategoryHandler) RegisterRoutes(router fiber.Router) {
	categories := router.Group("/categories")
	categories.Post("/", middlewares.RequireAuth(), middlewares.Authorize("categories", "write"), h.Create)
	categories.Get("/:id", h.GetByID)
	categories.Get("/", h.GetAll)
	categories.Delete("/:id", middlewares.RequireAuth(), middlewares.Authorize("categories", "write"), h.Delete)
	categories.Put("/:id", middlewares.RequireAuth(), middlewares.Authorize("categories", "write"), h.Update)
	categories.Get("/by-slug/:slug", h.GetBySlug)

	// categories.Get("/filters/:category_id", h.GetFiltersByCategoryID)
}
