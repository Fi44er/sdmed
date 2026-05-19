package user_http

import (
	"github.com/Fi44er/sdmed/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users", middlewares.RequireAuth())

	users.Get("/", middlewares.Authorize("users", "all"), h.GetAll)
	users.Get("/me", h.GetMy)
	users.Get("/:id", middlewares.Authorize("users", "all"), h.GetByID)
	users.Post("/", middlewares.Authorize("users", "all"), h.Create)
	users.Put("/:id", middlewares.Authorize("users", "all"), h.Update)
	users.Delete("/:id", middlewares.Authorize("users", "all"), h.Delete)
}
