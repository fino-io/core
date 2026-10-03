package core

import (
	"errors"
	"fmt"
	"net/http"
)

func NewError(code *ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewErrorf(code *ErrorCode, format string, arguments ...any) *Error {
	return NewError(code, fmt.Sprintf(format, arguments...))
}

func NewErrorFrom(code int32, message string) *Error {
	return NewError(NewErrorCode(code), message)
}

func NewFormattedErrorFrom(code int32, format string, arguments ...any) *Error {
	return NewErrorFrom(code, fmt.Sprintf(format, arguments...))
}

func (e *Error) Is(target error) bool {
	_, ok := target.(*Error)
	return ok
}

func IsError(err error) bool {
	return AsError(err) != nil
}

func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if len(e.Message) == 0 {
		if e.Code == nil {
			return ""
		}
		return e.Code.Name
	}
	return e.Message
}

func (e *Error) StatusCode() int {
	if e == nil || e.Code == nil {
		return http.StatusInternalServerError
	}
	if e.Code.HttpStatusCode >= http.StatusContinue && e.Code.HttpStatusCode <= 599 {
		return int(e.Code.HttpStatusCode)
	}
	return http.StatusInternalServerError
}

func (e *Error) AddDetail(detail any) *Error {
	if e != nil {
		v, _ := NewValue(detail)
		e.Details = append(e.Details, v)
	}
	return e
}
