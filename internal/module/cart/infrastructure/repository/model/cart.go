package cart_model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Cart struct {
	ID         string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID     string     `gorm:"type:uuid;index"`
	TotalPrice float64    `gorm:"type:decimal(10,2)"`
	Items      []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type CartItem struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CartID          string         `gorm:"type:uuid;index"`
	ProductID       string         `gorm:"type:uuid;index"`
	Article         string         `gorm:"type:text"`
	Quantity        int            `gorm:"type:integer"`
	UnitPrice       float64        `gorm:"type:decimal(10,2)"`
	SelectedOptions datatypes.JSON `gorm:"type:jsonb"`
	Iso             string         `gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
