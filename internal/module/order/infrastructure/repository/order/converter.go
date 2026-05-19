package order_repository

import (
	"encoding/json"

	order_entity "github.com/Fi44er/sdmed/internal/module/order/entity"
	order_model "github.com/Fi44er/sdmed/internal/module/order/infrastructure/repository/model"
	"gorm.io/datatypes"
)

type Converter struct{}

func (c *Converter) ToEntity(m *order_model.Order) *order_entity.Order {
	items := make([]order_entity.OrderItem, len(m.Items))
	for i, item := range m.Items {
		items[i] = *c.ToItemEntity(&item)
	}

	return &order_entity.Order{
		ID:           m.ID,
		UserID:       m.UserID,
		TotalPrice:   m.TotalPrice,
		Email:        m.Email,
		Phone:        m.Phone,
		FIO:          m.FIO,
		Address:      m.Address,
		Status:       order_entity.OrderStatus(m.Status),
		FragmentLink: m.FragmentLink,
		Items:        items,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (c *Converter) ToModel(e *order_entity.Order) *order_model.Order {
	return &order_model.Order{
		ID:           e.ID,
		UserID:       e.UserID,
		TotalPrice:   e.TotalPrice,
		Email:        e.Email,
		Phone:        e.Phone,
		FIO:          e.FIO,
		Address:      e.Address,
		Status:       string(e.Status),
		FragmentLink: e.FragmentLink,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func (c *Converter) ToItemEntity(m *order_model.OrderItem) *order_entity.OrderItem {
	var options any
	_ = json.Unmarshal(m.SelectedOptions, &options)

	return &order_entity.OrderItem{
		ID:              m.ID,
		OrderID:         m.OrderID,
		ProductID:       m.ProductID,
		Name:            m.Name,
		Article:         m.Article,
		Price:           m.Price,
		Quantity:        m.Quantity,
		TotalPrice:      m.TotalPrice,
		SelectedOptions: options,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func (c *Converter) ToItemModel(e *order_entity.OrderItem) *order_model.OrderItem {
	options, _ := json.Marshal(e.SelectedOptions)

	return &order_model.OrderItem{
		ID:              e.ID,
		OrderID:         e.OrderID,
		ProductID:       e.ProductID,
		Name:            e.Name,
		Article:         e.Article,
		Price:           e.Price,
		Quantity:        e.Quantity,
		TotalPrice:      e.TotalPrice,
		SelectedOptions: datatypes.JSON(options),
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
