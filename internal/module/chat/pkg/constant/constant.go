package chat_constant

import "github.com/Fi44er/sdmed/pkg/customerr"

var (
	ErrInvalidRole = customerr.NewError(400, "invalid role")
)
