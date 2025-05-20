package http

import "github.com/gofiber/fiber/v2"

func (h *UserHandler) RegisterRoutes(router fiber.Router) {
	users := router.Group("/users")

	users.Get("/", h.GetAll)
	users.Get("/me", h.GetMy)
	users.Get("/:id", h.GetByID)
	users.Post("/", h.Create)
	users.Put("/:id", h.Update)
	users.Delete("/:id", h.Delete)
}
