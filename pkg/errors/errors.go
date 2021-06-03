package errors

import (
	"net/http"
)

// ErrorCode represents different application error codes.
type ErrorCode string

// Enumerated values for ErrorCode
const (
	ServerErrorCode       = ErrorCode("server_error")
	BadRequestErrorCode   = ErrorCode("bad_request")
	UnauthorizedErrorCode = ErrorCode("unauthorized")
	NotFoundErrorCode     = ErrorCode("not_found")
	ConflictErrorCode     = ErrorCode("conflict")
)

// Error respresents an application error
type Error struct {
	Err        error
	StatusCode int
	Code       ErrorCode
	Message    string
}

func (e *Error) Error() string {
	var errMsg string
	if e.Message != "" {
		errMsg = e.Message
	}
	if e.Err != nil {
		if errMsg != "" {
			errMsg += ": "
		}
		errMsg = e.Err.Error()
	}
	return errMsg
}

// WithWrap wraps native error into application error
func (e *Error) WithWrap(err error) *Error {
	e.Err = err
	return e
}

// NewInternalServerError create a new application error with status 500
func NewInternalServerError(err error) *Error {
	return &Error{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		Code:       ServerErrorCode,
		Message:    http.StatusText(http.StatusInternalServerError),
	}
}

// NewBadRequestError create a new application error with status 400
func NewBadRequestError(message string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       BadRequestErrorCode,
		Message:    message,
	}
}

// NewUnauthorizedError create a new application error with status 401
func NewUnauthorizedError() *Error {
	return &Error{
		StatusCode: http.StatusUnauthorized,
		Code:       UnauthorizedErrorCode,
		Message:    http.StatusText(http.StatusUnauthorized),
	}
}

// NewNotFoundError create a new application error with status 404
func NewNotFoundError() *Error {
	return &Error{
		StatusCode: http.StatusNotFound,
		Code:       NotFoundErrorCode,
		Message:    http.StatusText(http.StatusNotFound),
	}
}

// NewConflictError create a new application error with status 409
func NewConflictError(message string) *Error {
	return &Error{
		StatusCode: http.StatusConflict,
		Code:       "conflict",
		Message:    message,
	}
}
