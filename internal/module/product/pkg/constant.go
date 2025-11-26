package product_constant

import "github.com/Fi44er/sdmed/pkg/customerr"

var (
	ErrCategoryNotFound     = customerr.NewError(404, "category not found")
	ErrCategoryAlreadyExist = customerr.NewError(409, "category already exist")

	ErrInvalidDataTypeCharacteristic = customerr.NewError(400, "invalid data type for characteristic")
)
