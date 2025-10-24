package contracts

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/product/entity"
)

type IFileUsecase interface {
	UploadFile(ctx context.Context, file *product_entity.File) error
}

type IProductRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]product_entity.Product, error)
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
	Create(ctx context.Context, entity *product_entity.Product) error
	Update(ctx context.Context, entity *product_entity.Product) error
	Delete(ctx context.Context, id string) error
}
