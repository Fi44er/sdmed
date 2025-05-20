package customerr

import "errors"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *Error) Is(target error) bool {
	if err, ok := target.(*Error); ok {
		return err.Code == e.Code && err.Message == e.Message
	}
	return errors.Is(e.Cause, target)
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func WrapError(code int, message string, cause error) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}
