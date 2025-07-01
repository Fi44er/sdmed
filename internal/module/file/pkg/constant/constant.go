package constant

import "github.com/Fi44er/sdmed/pkg/customerr"

var (
	ErrFileNotFound = customerr.NewError(404, "File not found")
)
