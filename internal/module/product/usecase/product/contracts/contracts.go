package product_usecase_contracts

import (
	"context"
	"time"

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
	GetAll(ctx context.Context, params product_entity.ProductFilterParams) ([]product_entity.Product, int64, error)
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
	Create(ctx context.Context, entity *product_entity.Product) error
	Update(ctx context.Context, entity *product_entity.Product) error
	Delete(ctx context.Context, id string) error
	GetByArticle(ctx context.Context, article string) (*product_entity.Product, error)
	GetBySlug(ctx context.Context, slug string) (*product_entity.Product, error)
	Count(ctx context.Context) (int64, error)
	GetFiltersByCategory(ctx context.Context, categoryID string) ([]product_entity.Filter, error)
}

type ICache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
}

type ICharValueUsecase interface {
	CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error
}
