package characteristic_repository

import (
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *product_entity.Characteristic) *product_model.Characteristic {
	options := make([]product_model.CharOption, len(entity.Options))
	for i, option := range entity.Options {
		options[i] = product_model.CharOption{
			ID:    option.ID,
			Value: option.Value,
		}
	}

	return &product_model.Characteristic{
		ID:          entity.ID,
		Name:        entity.Name,
		CategoryID:  entity.CategoryID,
		Unit:        entity.Unit,
		Description: entity.Description,
		DataType:    product_model.DataType(entity.DataType),
		Options:     options,
		IsRequired:  entity.IsRequired,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

func (c *Converter) ToEntity(model *product_model.Characteristic) *product_entity.Characteristic {
	return &product_entity.Characteristic{
		ID:          model.ID,
		Name:        model.Name,
		CategoryID:  model.CategoryID,
		Unit:        model.Unit,
		Description: model.Description,
		DataType:    product_entity.DataType(model.DataType),
		IsRequired:  model.IsRequired,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
