package char_value_repository

import (
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *product_entity.ProductCharValue) *product_model.CharacteristicValue {
	return &product_model.CharacteristicValue{
		CharacteristicID: entity.CharacteristicID,
		ProductID:        entity.ProductID,
		StringValue:      *entity.StringValue,
		NumberValue:      *entity.NumberValue,
		BooleanValue:     *entity.BooleanValue,
		OptionID:         *entity.OptionID,
	}
}
