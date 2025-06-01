package middleware

import (
	customa_err "github.com/Fi44er/sdmed/pkg/customerr"
	"github.com/gofiber/fiber/v2"
)

func ErrHandler(ctx *fiber.Ctx) error {
	err := ctx.Next()

	if err != nil {
		if error, ok := err.(*customa_err.Error); ok {
			return ctx.Status(error.Code).JSON(fiber.Map{
				"status":  "failed",
				"message": error.Message,
			})
		}
		return ctx.Status(500).JSON(fiber.Map{
			"status":  "failed",
			"message": "internal server error",
		})
	}

	return nil
}
