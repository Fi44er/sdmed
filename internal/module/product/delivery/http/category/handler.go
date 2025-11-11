package category_http

import (
	"context"

	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
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

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body product_dto.CreateCategoryDTO true "Category"
// @Success 201 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /categories [post]
func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	dto := new(product_dto.CreateCategoryDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntity, h.logger)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := h.usecase.Create(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error creating category: %v", err)
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "category created successfully",
	})
}
