package cart_dto

import cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"

type AddCartItemRequest struct {
	ProductID      string                      `json:"product_id" validate:"required"`
	Quantity       int                         `json:"quantity" validate:"required,min=1"`
	Iso            string                      `json:"iso"`
	DynamicOptions []cart_entity.DynamicOption `json:"dynamic_options"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required"`
}

type CartItemResponse struct {
	ID              string                               `json:"id"`
	ProductID       string                               `json:"product_id"`
	Article         string                               `json:"article"`
	Name            string                               `json:"name"`
	Image           string                               `json:"image"`
	Quantity        int                                  `json:"quantity"`
	UnitPrice       float64                              `json:"unit_price"`
	TotalPrice      float64                              `json:"total_price"`
	SelectedOptions []cart_entity.SelectedOptionResponse `json:"selected_options"`
	Iso             string                               `json:"iso"`
}

type CartResponse struct {
	ID         string             `json:"id"`
	Items      []CartItemResponse `json:"items"`
	TotalQty   int                `json:"total_qty"`
	TotalPrice float64            `json:"total_price"`
}

type MoveCartRequest struct {
	UserID string `json:"user_id"`
}
