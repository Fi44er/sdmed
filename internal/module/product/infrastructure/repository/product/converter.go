package product_repository

import (
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *product_entity.Product) *product_model.Product {
	return &product_model.Product{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price.Price,
		CategoryID:  entity.CategoryID,
	}
}

func (c *Converter) ToEntity(model *product_model.Product) *product_entity.Product {
	return &product_entity.Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Price:       product_entity.Price{Price: model.Price},
		CategoryID:  model.CategoryID,
	}
}
