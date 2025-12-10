package product_model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Slug      string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name      string         `gorm:"type:varchar(255);not null;uniqueIndex"`
	CreatedAt time.Time      `gorm:"not null;default:now()"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Characteristics []Characteristic `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}
