package product_model

import "time"

type Product struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Article     string `gorm:"type:varchar(255);unique;not null"`
	Name        string `gorm:"type:varchar(255);not null"`
	Slug        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:varchar(255);"`

	CategoryID      string                `gorm:"type:varchar(255);not null"`
	Characteristics []CharacteristicValue `gorm:"foreignKey:ProductID" json:"characteristics"`

	ManualPrice    float64 `gorm:"type:decimal(10,2);not null"`
	UseManualPrice bool    `gorm:"type:boolean;default:false"`

	IsActive bool `gorm:"type:boolean;default:true"`

	CreatedAt time.Time `gorm:"type:timestamp;default:now();"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:now();"`
}
