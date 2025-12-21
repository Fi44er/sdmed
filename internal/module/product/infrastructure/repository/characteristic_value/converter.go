package char_value_repository

import (
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *product_entity.ProductCharValue) *product_model.CharacteristicValue {
	model := &product_model.CharacteristicValue{
		CharacteristicID: entity.CharacteristicID,
		ProductID:        entity.ProductID,
		StringValue:      entity.StringValue,
		NumberValue:      entity.NumberValue,
		BooleanValue:     entity.BooleanValue,
	}

	if entity.OptionID != nil && *entity.OptionID != "" {
		model.OptionID = entity.OptionID
	} else {
		model.OptionID = nil
	}

	return model
}
