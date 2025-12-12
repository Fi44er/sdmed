package category_usecase_contracts

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

type ICategoryRepository interface {
	Create(ctx context.Context, entity *product_entity.Category) error
	Update(ctx context.Context, category *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
	GetBySlug(ctx context.Context, slug string) (*product_entity.Category, error)
	GetByName(ctx context.Context, name string) (*product_entity.Category, error)
	GetAll(ctx context.Context, offset, limit int) ([]product_entity.Category, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

type ICharacteristicUsecase interface {
	Create(ctx context.Context, characteristic *product_entity.Characteristic) error
	CreateMany(ctx context.Context, characteristics []product_entity.Characteristic) error
	Delete(ctx context.Context, id string) error
}
