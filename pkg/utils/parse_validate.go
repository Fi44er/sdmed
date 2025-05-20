package utils

import (
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ParseAndValidate[T any, D any](
	ctx *fiber.Ctx,
	dto *T,
	validator *validator.Validate,
	convert func(*T) *D,
	logger *logger.Logger,
) (*D, error) {

	// Парсим JSON в DTO
	if err := ctx.BodyParser(dto); err != nil {
		logger.Warnf("error while parsing body: %s", err)
		return nil, ctx.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Валидируем DTO
	if err := validator.Struct(dto); err != nil {
		logger.Warnf("error while validating dto: %s", err)
		return nil, ctx.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	// Конвертируем DTO → Domain
	domain := convert(dto)
	return domain, nil
}
