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
	GetAll(ctx context.Context, params *product_entity.ProductFilterParams) ([]product_entity.Product, int64, error)
	GetFilters(ctx context.Context, categoryID string) ([]product_entity.Filter, error)
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

// @Summary Get all products with filters
// @Description Get a list of products with support for pagination, sorting, and dynamic filters by characteristics
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Items per page (default 10)"
// @Param category_id query string false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param sort query string false "Sorting order: price_asc, price_desc, newest"
// @Param chars query []string false "Dynamic filters in format chars[char_id]=value"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /products [get]
func (h *ProductHandler) GetAll(ctx *fiber.Ctx) error {
	params := &product_dto.ProductQueryParams{
		Page:            1,
		PageSize:        10,
		Characteristics: make(map[string]string),
	}

	if err := ctx.QueryParser(params); err != nil {
		h.logger.Errorf("Failed to parse query params: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
	}

	allQueries := ctx.Queries()
	for key, value := range allQueries {
		if len(key) > 7 && key[:6] == "chars[" && key[len(key)-1:] == "]" {
			charID := key[6 : len(key)-1]
			if value != "" {
				params.Characteristics[charID] = value
			}
		}
	}

	h.logger.Debugf("Parsed query params: %+v", params)

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	filterParamsEntity := h.converter.ToFilterEntity(*params)

	products, count, err := h.usecase.GetAll(ctx.Context(), &filterParamsEntity)
	if err != nil {
		return err
	}

	productRes := h.converter.ToProductListResponse(products, count, params.Page, params.PageSize)

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   productRes,
	})
}

// @Summary Get filters for a category
// @Description Get filters for a category
// @Tags products
// @Accept json
// @Produce json
// @Param category_id path string true "Category ID"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /products/filters/{category_id} [get]
func (h *ProductHandler) GetFilters(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("category_id")

	filters, err := h.usecase.GetFilters(ctx.Context(), categoryID)
	if err != nil {
		return err
	}

	filtersRes := h.converter.ToFilterResponses(filters)

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   filtersRes,
	})
}
