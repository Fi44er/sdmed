package customerr

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e.Code == t.Code && e.Message == t.Message
	}
	return false
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) WithCause(cause error) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Cause:   cause,
	}
}

func (e *Error) WithContext(contextMsg string) *Error {
	return &Error{
		Code:    e.Code,
		Message: fmt.Sprintf("%s: %s", contextMsg, e.Message),
		Cause:   e.Cause,
	}
}

func FromError(err error) (int, string) {
	if err == nil {
		return 200, ""
	}
	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr.Code, customErr.Message
	}
	return 500, "Internal Server Error"
}
