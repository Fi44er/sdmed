package middlewares

import (
	accessmanager_service "github.com/Fi44er/sdmed/internal/module/auth/usecase/access_manager"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/gofiber/fiber/v2"
)

const ManagerKey = "accessManager"

func InjectManager(am *accessmanager_service.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(ManagerKey, am)
		return c.Next()
	}
}

func Authorize(obj, act string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		val := ctx.Locals(ManagerKey)
		if val == nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Access manager not found",
			})
		}
		am := val.(*accessmanager_service.Manager)

		sess := session.FromFiberContext(ctx)
		if sess == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		return ctx.Next()
	}
}
