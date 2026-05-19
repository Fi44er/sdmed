package order_dto

import (
	order_entity "github.com/Fi44er/sdmed/internal/module/order/entity"
)

type CreateOrderRequest struct {
	FIO         string `json:"fio" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Address     string `json:"address" validate:"required"`
}

type CreateDirectOrderRequest struct {
	CreateOrderRequest
	ProductID string `json:"product_id" validate:"required,uuid4"`
}

type ChangeOrderStatusRequest struct {
	OrderID string                   `json:"order_id" validate:"required,uuid4"`
	Status  order_entity.OrderStatus `json:"status" validate:"required"`
}
