package cart_adapters

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/product"
)

type ProductUsecaseAdapter struct {
	productUsecase product_usecase.IProductUsecase
}

func NewProductUsecaseAdapter(productUsecase product_usecase.IProductUsecase) *ProductUsecaseAdapter {
	return &ProductUsecaseAdapter{
		productUsecase: productUsecase,
	}
}

func (a *ProductUsecaseAdapter) GetByID(ctx context.Context, id string) (*product_entity.Product, error) {
	return a.productUsecase.GetByID(ctx, id)
}
