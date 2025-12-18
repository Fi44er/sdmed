package product_model

import "time"

type CharacteristicValue struct {
	ID               string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	CharacteristicID string `gorm:"type:varchar(255);not null"`
	ProductID        string `gorm:"type:varchar(255);not null"`

	StringValue  string  `gorm:"type:varchar(255);"`
	NumberValue  float64 `gorm:"type:float;"`
	BooleanValue bool    `gorm:"type:boolean;"`

	OptionID string `gorm:"type:varchar(255);"`

	Option  CharOption `gorm:"foreignKey:OptionID;references:ID"`
	Product Product    `gorm:"foreignKey:ProductID;references:ID"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
