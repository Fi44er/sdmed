package chat_constant

import "github.com/Fi44er/sdmed/pkg/customerr"

var (
	ErrInvalidRole     = customerr.NewError(400, "invalid role")
	ErrChatNotFound    = customerr.NewError(404, "chat not found")
	ErrMessageNotFound = customerr.NewError(404, "message not found")
	ErrChatClosed      = customerr.NewError(400, "chat is closed")
	ErrUnauthorized    = customerr.NewError(403, "unauthorized to access this chat")
)
