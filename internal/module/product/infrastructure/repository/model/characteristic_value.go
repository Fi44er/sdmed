package product_model

import "time"

type CharacteristicValue struct {
	ID               string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	CharacteristicID string `gorm:"type:varchar(255);not null"`
	ProductID        string `gorm:"type:varchar(255);not null"`

	StringValue  string  `gorm:"type:varchar(255);not null"`
	NumberValue  float64 `gorm:"type:float;not null"`
	BooleanValue bool    `gorm:"type:boolean;not null"`

	OptionID string `gorm:"type:varchar(255);not null"`

	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}
