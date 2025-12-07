package category_http

import (
	"context"

	"github.com/Fi44er/sdmed/internal/config"
	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/Fi44er/sdmed/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ICategoryUsecase interface {
	Create(ctx context.Context, category *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, category *product_entity.Category) error
}

type CategoryHandler struct {
	usecase ICategoryUsecase

	validator *validator.Validate
	logger    *logger.Logger
	converter *Converter
	config    *config.Config
}

func NewCategoryHandler(
	usecase ICategoryUsecase,
	logger *logger.Logger,
	validator *validator.Validate,
	config *config.Config,
) *CategoryHandler {
	return &CategoryHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		config:    config,
		converter: NewConverter(config),
	}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body product_dto.CreateCategoryRequest true "Category"
// @Success 201 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /categories [post]
func (h *CategoryHandler) Create(ctx *fiber.Ctx) error {
	dto := new(product_dto.CreateCategoryRequest)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityFromCreate, h.logger)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := h.usecase.Create(ctx.Context(), entity); err != nil {
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "category created successfully",
	})
}

// Update godoc
// @Summary Update a category
// @Description Update a category
// @Description ### Особенности обновления характеристик:
// @Description **Важно:** При обновлении категории поле `Characteristics` работает по принципу полной замены.
// @Description 1. **Добавление новых характеристик:**
// @Description   Укажите все существующие характеристики + новые
// @Description 2. **Удаление характеристик:**
// @Description   Укажите только те характеристики, которые должны остаться
// @Description 3. **Изменение характеристик:**
// @Description   Укажите обновленный список всех характеристик
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body product_dto.UpdateCategoryRequest true "Category"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /categories/{id} [put]
func (h *CategoryHandler) Update(ctx *fiber.Ctx) error {
	dto := new(product_dto.UpdateCategoryRequest)
	dto.ID = ctx.Params("id")

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityFromUpdate, h.logger)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := h.usecase.Update(ctx.Context(), entity); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "category updated successfully",
	})
}

// GetByID godoc
// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.ResponseData{data=product_dto.CategoryResponse} "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	category, err := h.usecase.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	categoryRes := h.converter.toCategoryResponse(category)

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   categoryRes,
	})
}

// GetAll godoc
// @Summary Get all categories
// @Description Get list of categories with pagination
// @Tags categories
// @Accept json
// @Produce json
// @Param offset path int false "Offset for pagination" default(0)
// @Param limit path int false "Limit for pagination" default(10)
// @Success 200 {object} response.ResponseData{data=[]product_dto.CategoryResponse} "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /categories [get]
func (h *CategoryHandler) GetAll(ctx *fiber.Ctx) error {
	offset := ctx.QueryInt("offset")
	limit := ctx.QueryInt("limit")

	categories, err := h.usecase.GetAll(ctx.Context(), offset, limit)
	if err != nil {
		return err
	}
	categoriesRes := h.converter.toCategoryResponses(categories)

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   categoriesRes,
	})
}

// Delete godoc
// @Summary Delete category
// @Description Delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.usecase.Delete(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "category deleted successfully",
	})
}
