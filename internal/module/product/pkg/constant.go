package product_constant

import "github.com/Fi44er/sdmed/pkg/customerr"

var (
	ErrCategoryNotFound      = customerr.NewError(404, "category not found")
	ErrCategoryAlreadyExists = customerr.NewError(409, "category already exist")

	ErrInvalidDataTypeCharacteristic = customerr.NewError(400, "invalid data type for characteristic")
	ErrCharacteristicAlreadyExists   = customerr.NewError(409, "characteristic already exists")
)
