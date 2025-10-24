package category_repository

import (
	"github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/internal/module/product/infrastucture/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(product_entity *product_entity.Category) *product_model.Category {
	return &product_model.Category{
		ID:   product_entity.ID,
		Name: product_entity.Name,
	}
}

func (c *Converter) Toproduct_entity(model *product_model.Category) *product_entity.Category {
	return &product_entity.Category{
		ID:   model.ID,
		Name: model.Name,
	}
}
