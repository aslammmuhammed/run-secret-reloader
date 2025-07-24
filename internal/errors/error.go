package error

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code     string
	Message  string
	HTTPCode int
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewNotFoundError creates a new AppError for 404 Not Found responses
// Example usage:
//
//	return appError.NewNotFoundError("user not found")
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:     http.StatusText(http.StatusNotFound),
		HTTPCode: http.StatusNotFound,
		Message:  message,
	}
}

// NewBadRequestError creates a new AppError for 400 Bad Request responses
// Example usage:
//
//	return appError.NewBadRequestError("invalid authentication token")
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:     http.StatusText(http.StatusBadRequest),
		HTTPCode: http.StatusBadRequest,
		Message:  message,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:     http.StatusText(http.StatusInternalServerError),
		HTTPCode: http.StatusInternalServerError,
		Message:  message,
	}
}

func NewAuthenticationError(message string) *AppError {
	return &AppError{
		Code:     http.StatusText(http.StatusUnauthorized),
		HTTPCode: http.StatusUnauthorized,
		Message:  message,
	}
}

func NewAlreadyExistsError(message string) *AppError {
	return &AppError{
		Code:     http.StatusText(http.StatusConflict),
		HTTPCode: http.StatusConflict,
		Message:  message,
	}
}

func NewAuthorizationError(message string) *AppError {
	return &AppError{
		Code:     http.StatusText(http.StatusForbidden),
		HTTPCode: http.StatusForbidden,
		Message:  message,
	}
}
