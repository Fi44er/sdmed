package middlewares

import (
	"fmt"

	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	accessmanager_service "github.com/Fi44er/sdmed/internal/module/auth/usecase/access_manager"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gofiber/fiber/v2"
)

func InjectManager(am *accessmanager_service.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(ManagerKey, am)
		return c.Next()
	}
}

const ManagerKey = "accessManager"

func Guest() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess := session.FromFiberContext(ctx)
		if sess == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		sessionData, ok := sess.Get("session_info").(map[string]any)
		if !ok || sessionData == nil {
			data := auth_entity.ActiveSession{
				UserRoles: []string{"guest"},
				IsShadow:  true,
			}

			var newData map[string]any
			if err := mapstructure.Decode(data, &newData); err != nil {
				return fmt.Errorf("failed to encode session data for guest user: %w", err)
			}
			sess.Put("session_info", newData)
		}

		return ctx.Next()
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

		sessionData, ok := sess.Get("session_info").(map[string]any)
		if !ok {
			return ctx.Status(403).JSON(fiber.Map{
				"error": "Invalid session data",
			})
		}

		var userSession auth_entity.ActiveSession
		if err := mapstructure.Decode(sessionData, &userSession); err != nil {
			return fmt.Errorf("failed to decode session data: %v", err)
		}

		hasAccess := false
		if userSession.UserRoles != nil {
			for _, role := range userSession.UserRoles {
				ok, _ := am.Enforcer.Enforce(role, obj, act)
				if ok {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden",
			})
		}

		return ctx.Next()
	}
}

func RequireAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess := session.FromFiberContext(ctx)
		if sess == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session not found",
			})
		}

		sessionData, ok := sess.Get("session_info").(map[string]any)
		if !ok || sessionData == nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Not authenticated",
			})
		}

		var userSession auth_entity.ActiveSession
		if err := mapstructure.Decode(sessionData, &userSession); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid session data",
			})
		}

		// Проверяем, что пользователь не shadow
		if userSession.IsShadow {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Authentication required",
				"message": "Please sign in to access this resource",
			})
		}

		return ctx.Next()
	}
}
