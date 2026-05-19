package cart_http

import (
	cart_dto "github.com/Fi44er/sdmed/internal/module/cart/dto"
	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
)

type Converter struct{}

func (c *Converter) ToCartResponse(cart *cart_entity.Cart) *cart_dto.CartResponse {
	items := make([]cart_dto.CartItemResponse, 0, len(cart.Items))
	totalQty := 0

	for _, item := range cart.Items {
		itemTotalPrice := item.UnitPrice * float64(item.Quantity)
		items = append(items, cart_dto.CartItemResponse{
			ID:              item.ID,
			ProductID:       item.ProductID,
			Article:         item.Article,
			Quantity:        item.Quantity,
			UnitPrice:       item.UnitPrice,
			TotalPrice:      itemTotalPrice,
			Iso:             item.Iso,
			SelectedOptions: c.toSelectedOptions(item.SelectedOptions),
		})
		totalQty += item.Quantity
	}

	return &cart_dto.CartResponse{
		ID:         cart.ID,
		Items:      items,
		TotalQty:   totalQty,
		TotalPrice: cart.TotalPrice,
	}
}

func (c *Converter) toSelectedOptions(raw any) []cart_entity.SelectedOptionResponse {
	if raw == nil {
		return []cart_entity.SelectedOptionResponse{}
	}

	if opts, ok := raw.([]cart_entity.SelectedOptionResponse); ok {
		return opts
	}

	// Fallback/Parsing if needed
	return []cart_entity.SelectedOptionResponse{}
}
