package product_http

import (
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IProductUsecase interface{}

type ProductHandler struct {
	usecase IProductUsecase

	validator *validator.Validate
	logger    *logger.Logger
}

func NewProductHandler(
	usecase IProductUsecase,
	validator *validator.Validate,
	logger *logger.Logger,
) *ProductHandler {
	return &ProductHandler{
		usecase:   usecase,
		validator: validator,
		logger:    logger,
	}
}

func (h *ProductHandler) Create(ctx *fiber.Ctx) error {

	return ctx.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "product created successfully",
	})
}
