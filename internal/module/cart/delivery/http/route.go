package cart_http

import "github.com/gofiber/fiber/v2"

func (h *CartHandler) RegisterRoutes(router fiber.Router) {
	cart := router.Group("/cart")
	
	cart.Get("/", h.Get)
	cart.Post("/items", h.AddItem)
	cart.Delete("/items/:id", h.DeleteItem)
	cart.Post("/move", h.Move)
}
