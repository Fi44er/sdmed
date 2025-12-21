package product_entity

import (
	"strconv"
	"time"
)

type ProductCharValue struct {
	ID               string
	CharacteristicID string
	ProductID        string

	StringValue  *string
	NumberValue  *float64
	BooleanValue *bool

	OptionID *string
	Option   *CharOption

	CharacteristicName string

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (cv *ProductCharValue) GetStringValue() string {
	if cv.StringValue != nil && *cv.StringValue != "" {
		return *cv.StringValue
	}
	if cv.Option != nil && cv.Option.Value != "" {
		return cv.Option.Value
	}
	if cv.NumberValue != nil {
		return strconv.FormatFloat(*cv.NumberValue, 'f', -1, 64)
	}
	if cv.BooleanValue != nil {
		return strconv.FormatBool(*cv.BooleanValue)
	}
	return ""
}
