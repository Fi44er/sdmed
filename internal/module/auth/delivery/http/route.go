package auth_http

import (
	"github.com/Fi44er/sdmed/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/sign-up", h.SignUp)
	auth.Post("/sign-in", h.SignIn)
	auth.Post("/verify-code", h.VerifyCode)
	auth.Post("/send-code", h.SendCode)
	auth.Post("/forgot-password", h.ForgotPassword)
	auth.Get("/validate-reset-password", h.ValidateResetPassword)
	auth.Post("/reset-password", h.ResetPassword)

	protected := auth.Group("", middlewares.RequireAuth())
	{
		// Обновление сессии (аналог refresh token)
		protected.Post("/refresh", h.RefreshSession)

		// Выход
		protected.Post("/sign-out", h.SignOut)
		protected.Post("/sign-out-all", h.SignOutAll)

		// Управление устройствами
		protected.Get("/devices", h.GetDevices)
		protected.Post("/devices/:device_id/revoke", h.RevokeDevice)
	}
}
