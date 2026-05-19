package order_entity

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunded  OrderStatus = "refunded"
	OrderStatusShpping   OrderStatus = "shipping"
	OrderStatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID           string
	UserID       *string
	TotalPrice   float64
	Email        string
	Phone        string
	FIO          string
	Address      string
	Status       OrderStatus
	FragmentLink string
	Items        []OrderItem
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OrderItem struct {
	ID              string
	OrderID         string
	ProductID       string
	Name            string
	Article         string
	Price           float64
	Quantity        int
	TotalPrice      float64
	SelectedOptions any // JSON data
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
