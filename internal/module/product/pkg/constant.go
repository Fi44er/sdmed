package product_constant

import "github.com/Fi44er/sdmed/pkg/customerr"

var (
	ErrCategoryNotFound = customerr.NewError(404, "category not found")
)
