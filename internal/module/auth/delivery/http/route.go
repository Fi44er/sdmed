package auth_http

import "github.com/gofiber/fiber/v2"

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/sign-in", h.SignIn)
	auth.Post("/verify-code", h.VerifyCode)
	auth.Post("/sign-out", h.SignOut)
	auth.Post("/refresh-token", h.RefreshToken)
	auth.Post("/send-code", h.SendCode)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Get("/validate-reset-password", h.ValidateResetPassword)
	auth.Post("/reset-password", h.ResetPassword)
}
