package order_model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Order struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID       *string        `gorm:"type:uuid;index"`
	TotalPrice   float64        `gorm:"type:decimal(10,2)"`
	Email        string         `gorm:"type:text"`
	Phone        string         `gorm:"type:text"`
	FIO          string         `gorm:"type:text"`
	Address      string         `gorm:"type:text"`
	Status       string         `gorm:"type:text;default:'pending'"`
	FragmentLink string         `gorm:"type:text"`
	Items        []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type OrderItem struct {
	ID              string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	OrderID         string         `gorm:"type:uuid;index"`
	ProductID       string         `gorm:"type:uuid;index"`
	Name            string         `gorm:"type:text"`
	Article         string         `gorm:"type:text"`
	Price           float64        `gorm:"type:decimal(10,2)"`
	Quantity        int            `gorm:"type:integer"`
	TotalPrice      float64        `gorm:"type:decimal(10,2)"`
	SelectedOptions datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
