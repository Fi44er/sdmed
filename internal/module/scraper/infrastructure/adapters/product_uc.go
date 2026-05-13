package scraper_adapters

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
)

type IProductUseCase interface {
	CreateMany(ctx context.Context, products []*product_entity.Product) error
}

type ProductUsecaseAdapter struct {
	productUsecase IProductUseCase
}

func NewProductUsecaseAdapter(productUsecase IProductUseCase) *ProductUsecaseAdapter {
	return &ProductUsecaseAdapter{productUsecase: productUsecase}
}

func (a *ProductUsecaseAdapter) CreateMany(ctx context.Context, products []*scraper_entity.Product) error {
	productEntities := make([]*product_entity.Product, 0, len(products))
	for _, product := range products {
		productEntities = append(productEntities, a.toProductEntity(product))
	}
	return a.productUsecase.CreateMany(ctx, productEntities)
}

func (a *ProductUsecaseAdapter) toProductEntity(product *scraper_entity.Product) *product_entity.Product {
	return &product_entity.Product{
		Article: product.Article,
		Name:    product.Name,
	}
}
