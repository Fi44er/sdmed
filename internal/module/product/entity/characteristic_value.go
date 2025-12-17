package product_entity

import "time"

type ProductCharValue struct {
	ID               string
	CharacteristicID string
	ProductID        string

	StringValue  *string
	NumberValue  *float64
	BooleanValue *bool

	OptionID *string
	// Option   *CharOption

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
