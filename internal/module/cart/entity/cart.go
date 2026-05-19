package cart_entity

import (
	"time"
)

type Cart struct {
	ID         string
	UserID     string
	Items      []CartItem
	TotalPrice float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CartItem struct {
	ID              string
	CartID          string
	ProductID       string
	Article         string
	Quantity        int
	UnitPrice       float64 // цена за единицу на момент добавления
	SelectedOptions any
	Iso             string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type DynamicOption struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type SelectedOptionResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
