package errors

import (
	"errors"
	"fmt"
)

// ErrorCode represents a unique identifier for each type of error
type ErrorCode string

const (
	// Common error codes
	ErrNotFound     ErrorCode = "NOT_FOUND"
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
	ErrInternal     ErrorCode = "INTERNAL_ERROR"
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"
	ErrConflict     ErrorCode = "CONFLICT"
	ErrBadRequest   ErrorCode = "BAD_REQUEST"
	ErrNotAvailable ErrorCode = "NOT_AVAILABLE"
)

// AppError represents a domain error in the application
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap implements the errors.Wrapper interface
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError without a wrapped error
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap creates a new AppError with a wrapped error
func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is implements error matching for AppError
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// GetCode extracts the ErrorCode from an error if it's an AppError
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if ok := As(err, &appErr); ok {
		return appErr.Code
	}
	return ErrInternal
}

// As provides error type assertion
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Is provides error type comparison
func Is(err, target error) bool {
	return errors.Is(err, target)
}
