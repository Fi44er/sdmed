package product_usecase_contracts

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type IFileUsecaseAdapter interface {
	MakeFilesPermanent(ctx context.Context, fileIDs []string, ownerID, ownerType string) error
	GetByOwner(ctx context.Context, ownerID, ownerType string) ([]product_entity.File, error)
	GetByOwners(ctx context.Context, ownerIDs []string, ownerType string) (map[string][]product_entity.File, error)

	DeleteByOwner(ctx context.Context, ownerID, ownerType string) error
	DeleteByID(ctx context.Context, id string) error
}

type IProductRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]product_entity.Product, error)
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
	Create(ctx context.Context, entity *product_entity.Product) error
	Update(ctx context.Context, entity *product_entity.Product) error
	Delete(ctx context.Context, id string) error
	GetByArticle(ctx context.Context, article string) (*product_entity.Product, error)
	GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error)
}

type ICharValueUsecase interface {
	CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error
}
