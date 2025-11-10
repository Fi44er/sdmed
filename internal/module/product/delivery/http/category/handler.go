package product_http

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ICategoryUsecase interface {
	Create(ctx context.Context, category *product_entity.Category) error
}

type CategoryHandler struct {
	usecase ICategoryUsecase

	validator *validator.Validate
	logger    *logger.Logger
	converter *Converter
}

func NewCategoryHandler(
	usecase ICategoryUsecase,
	logger *logger.Logger,
	validator *validator.Validate,
) *CategoryHandler {
	return &CategoryHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		converter: &Converter{},
	}
}

func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	return nil
}
