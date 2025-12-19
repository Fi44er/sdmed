package product_http

import (
	"context"

	"github.com/Fi44er/sdmed/internal/config"
	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IProductUsecase interface {
	Create(ctx context.Context, product *product_entity.Product) error
	GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error)
}

type ProductHandler struct {
	usecase IProductUsecase

	validator *validator.Validate
	logger    *logger.Logger
	converter *Converter
	config    *config.Config
}

func NewProductHandler(
	usecase IProductUsecase,
	validator *validator.Validate,
	logger *logger.Logger,
	config *config.Config,
) *ProductHandler {
	return &ProductHandler{
		usecase:   usecase,
		validator: validator,
		logger:    logger,
		converter: &Converter{config: config},
		config:    config,
	}
}

// Create godoc
// @Summary Create a product
// @Description Create a product
// @Tags products
// @Accept json
// @Produce json
// @Param category body product_dto.CreateProductRequest true "Product"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /products [post]
func (h *ProductHandler) Create(ctx *fiber.Ctx) error {
	dto := new(product_dto.CreateProductRequest)

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
		"message": "product created successfully",
	})
}

// @Summary Get a product by slug
// @Description Get a product by slug
// @Tags products
// @Accept json
// @Produce json
// @Param slug path string true "Slug"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /products/{slug} [get]
func (h *ProductHandler) GetBySlug(ctx *fiber.Ctx) error {
	slug := ctx.Params("slug")

	product, err := h.usecase.GetBySlug(ctx.Context(), slug)
	if err != nil {
		return err
	}

	productRes := h.converter.ToProductResponse(product)

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "product retrieved successfully",
		"data":    productRes,
	})
}
