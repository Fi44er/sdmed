package cart_repository

import (
	"encoding/json"

	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	cart_model "github.com/Fi44er/sdmed/internal/module/cart/infrastructure/repository/model"
	cart_repo_model "github.com/Fi44er/sdmed/internal/module/cart/infrastructure/repository/model"
	"gorm.io/datatypes"
)

type Converter struct{}

func (c *Converter) ToModel(entity *cart_entity.Cart) *cart_repo_model.Cart {
	items := make([]cart_model.CartItem, len(entity.Items))
	for i, item := range entity.Items {
		items[i] = *c.ToItemModel(&item)
	}
	return &cart_model.Cart{
		ID:         entity.ID,
		UserID:     entity.UserID,
		TotalPrice: entity.TotalPrice,
		Items:      items,
	}
}

func (c *Converter) ToEntity(model *cart_model.Cart) *cart_entity.Cart {
	items := make([]cart_entity.CartItem, len(model.Items))
	for i, item := range model.Items {
		items[i] = *c.ToItemEntity(&item)
	}
	return &cart_entity.Cart{
		ID:         model.ID,
		UserID:     model.UserID,
		TotalPrice: model.TotalPrice,
		Items:      items,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}

func (c *Converter) ToItemModel(entity *cart_entity.CartItem) *cart_model.CartItem {
	optionsJSON, _ := json.Marshal(entity.SelectedOptions)
	return &cart_model.CartItem{
		ID:              entity.ID,
		CartID:          entity.CartID,
		ProductID:       entity.ProductID,
		Article:         entity.Article,
		Quantity:        entity.Quantity,
		UnitPrice:       entity.UnitPrice,
		SelectedOptions: datatypes.JSON(optionsJSON),
		Iso:             entity.Iso,
	}
}

func (c *Converter) ToItemEntity(model *cart_model.CartItem) *cart_entity.CartItem {
	var options any
	_ = json.Unmarshal(model.SelectedOptions, &options)
	return &cart_entity.CartItem{
		ID:              model.ID,
		CartID:          model.CartID,
		ProductID:       model.ProductID,
		Article:         model.Article,
		Quantity:        model.Quantity,
		UnitPrice:       model.UnitPrice,
		SelectedOptions: options,
		Iso:             model.Iso,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}
