package model

import "fmt"

// EmptyResponse is an empty API response object
type EmptyResponse struct{}

// ECode is a error code (coarse-grained)
type ECode string

const (
	// EBadRequest is an error code for malformed arguments
	EBadRequest ECode = "bad_request"

	// ENotFound is an error code for a non-existing entity
	ENotFound ECode = "not_found"

	// EConflict is an error code for a conflicting entity
	EConflict ECode = "conflict"

	// EInternalError is an error code for an unexpected internal error
	EInternalError ECode = "internal_error"

	// EAccessDenied is an error code for an access error
	EAccessDenied ECode = "access_denied"
)

// Error is a service error object
type Error struct {
	Code    ECode  `json:"error"`
	Message string `json:"message"`
}

// Error implements error interface
func (e *Error) Error() string {
	return e.String()
}

// String converts object to a string
func (e *Error) String() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewError creates new error object
func NewError(code ECode, format string, a ...interface{}) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}
