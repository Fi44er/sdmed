package utils

import (
	"github.com/Fi44er/sdmed/pkg/logger"
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
	if err := ctx.BodyParser(dto); err != nil {
		logger.Warnf("error while parsing body: %s", err)
		return nil, err
	}

	if err := validator.Struct(dto); err != nil {
		logger.Warnf("error while validating dto: %s", err)
		return nil, err

	}

	domain := convert(dto)
	return domain, nil
}
