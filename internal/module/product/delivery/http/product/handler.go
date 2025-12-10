package product_http

import (
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
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
