package http

import "github.com/gofiber/fiber/v2"

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sign-up", h.SignUp)
}
