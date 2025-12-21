package product_model

import "time"

type CharacteristicValue struct {
	ID               string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	CharacteristicID string `gorm:"type:uuid;not null"`
	ProductID        string `gorm:"type:uuid;not null"`

	StringValue  *string  `gorm:"type:varchar(255);null"`
	NumberValue  *float64 `gorm:"type:float;null"`
	BooleanValue *bool    `gorm:"type:boolean;null"`

	OptionID *string `gorm:"type:uuid;null"`

	Option  CharOption `gorm:"foreignKey:OptionID;references:ID"`
	Product Product    `gorm:"foreignKey:ProductID;references:ID"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	CharacteristicName string `gorm:"->"`
}
