package char_value_usecase_contracts

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type ICharValueRepository interface {
	Create(ctx context.Context, charValue *product_entity.ProductCharValue) error
	Delete(ctx context.Context, id string) error
	CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error
}

type ICharacteristicUsecase interface {
	GetByIDs(ctx context.Context, ids []string) ([]product_entity.Characteristic, error)
}
