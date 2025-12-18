package char_value_repository

import (
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
)

type Converter struct{}

// func (c *Converter) ToModel(entity *product_entity.ProductCharValue) *product_model.CharacteristicValue {
// 	return &product_model.CharacteristicValue{
// 		CharacteristicID: entity.CharacteristicID,
// 		ProductID:        entity.ProductID,
// 		StringValue:      *entity.StringValue,
// 		NumberValue:      *entity.NumberValue,
// 		BooleanValue:     *entity.BooleanValue,
// 		OptionID:         *entity.OptionID,
// 	}
// }

func (c *Converter) ToModel(entity *product_entity.ProductCharValue) *product_model.CharacteristicValue {
	model := &product_model.CharacteristicValue{
		CharacteristicID: entity.CharacteristicID,
		ProductID:        entity.ProductID,
		// Инициализируем остальные поля дефолтными значениями
	}

	// Безопасное копирование StringValue
	if entity.StringValue != nil {
		model.StringValue = *entity.StringValue
	}

	// Безопасное копирование NumberValue
	if entity.NumberValue != nil {
		model.NumberValue = *entity.NumberValue
	} else {
		// Можно установить дефолтное значение или 0
		model.NumberValue = 0 // или float64(0)
	}

	// Безопасное копирование BooleanValue
	if entity.BooleanValue != nil {
		model.BooleanValue = *entity.BooleanValue
	} else {
		model.BooleanValue = false
	}

	// Безопасное копирование OptionID
	if entity.OptionID != nil {
		model.OptionID = *entity.OptionID
	}
	// Для OptionID можно оставить пустую строку, если nil

	return model
}
