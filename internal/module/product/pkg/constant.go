package product_constant

import (
	"time"

	"github.com/Fi44er/sdmed/pkg/customerr"
)

var (
	ErrCategoryNotFound      = customerr.NewError(404, "category not found")
	ErrCategoryAlreadyExists = customerr.NewError(409, "category already exist")

	ErrInvalidDataTypeCharacteristic = customerr.NewError(400, "invalid data type for characteristic")
	ErrCharacteristicAlreadyExists   = customerr.NewError(409, "characteristic already exists")
	ErrCharacteristicOptionsEmpty    = customerr.NewError(400, "characteristic options cannot be empty")

	ErrProductAlreadyExists = customerr.NewError(409, "product already exists")
	ErrProductNotFound      = customerr.NewError(404, "product not found")

	ErrCharacteristicNotFound      = customerr.NewError(404, "characteristic not found")
	ErrInvalidDataType             = customerr.NewError(400, "invalid data type")
	ErrValueRequired               = customerr.NewError(422, "value is required")
	ErrInvalidNumber               = customerr.NewError(400, "invalid number value")
	ErrInvalidBoolean              = customerr.NewError(400, "invalid boolean value")
	ErrOptionNotFound              = customerr.NewError(404, "option not found")
	ErrInvalidString               = customerr.NewError(400, "invalid string value")
	ErrRequiredCharacteristicEmpty = customerr.NewError(400, "required characteristic is empty")
	ErrInvalidValue                = customerr.NewError(400, "invalid value")
)

const (
	CategoryFiltersKeyPrefix = "filters:category:"
	FilterExpered            = time.Hour * 24
)
