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
		CategoryID:  entity.CategoryID,

		ManualPrice:    *entity.ManualPrice,
		UseManualPrice: entity.UseManualPrice,

		IsActive: entity.IsActive,
	}
}

func (c *Converter) ToEntity(model *product_model.Product) *product_entity.Product {
	return &product_entity.Product{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		CategoryID:  model.CategoryID,

		ManualPrice:    &model.ManualPrice,
		UseManualPrice: model.UseManualPrice,

		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
