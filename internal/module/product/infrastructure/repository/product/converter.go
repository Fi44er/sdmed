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
		Article:     entity.Article,
		Slug:        entity.Slug,
		Description: entity.Description,
		CategoryID:  entity.CategoryID,

		ManualPrice:    *entity.ManualPrice,
		UseManualPrice: entity.UseManualPrice,

		IsActive: entity.IsActive,
	}
}

func (c *Converter) ToEntity(model *product_model.Product) *product_entity.Product {
	charValues := make([]product_entity.ProductCharValue, 0)
	for _, charValue := range model.Characteristics {
		charValues = append(charValues, product_entity.ProductCharValue{
			ID:               charValue.ID,
			CharacteristicID: charValue.CharacteristicID,
			ProductID:        charValue.ProductID,
			StringValue:      &charValue.StringValue,
			NumberValue:      &charValue.NumberValue,
			BooleanValue:     &charValue.BooleanValue,
			OptionID:         &charValue.OptionID,
			Option:           (*product_entity.CharOption)(&charValue.Option),
			CreatedAt:        charValue.CreatedAt,
			UpdatedAt:        charValue.UpdatedAt,
		})
	}
	return &product_entity.Product{
		ID:          model.ID,
		Name:        model.Name,
		Article:     model.Article,
		Slug:        model.Slug,
		Description: model.Description,
		CategoryID:  model.CategoryID,
		CharValues:  charValues,

		ManualPrice:    &model.ManualPrice,
		UseManualPrice: model.UseManualPrice,

		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
