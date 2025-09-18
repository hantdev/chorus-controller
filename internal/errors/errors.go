package errors

import (
	"errors"
	"net/http"
)

// Error types for better error handling
var (
	ErrInvalidRequest      = errors.New("invalid request")
	ErrWorkerUnavailable   = errors.New("worker service unavailable")
	ErrReplicationNotFound = errors.New("replication not found")
	ErrInvalidReplication  = errors.New("invalid replication configuration")
)

// APIError represents an API error with HTTP status code
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// NewAPIError creates a new API error
func NewAPIError(code int, message string, err error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Common error constructors
func NewBadRequestError(message string, err error) *APIError {
	return NewAPIError(http.StatusBadRequest, message, err)
}

func NewInternalServerError(message string, err error) *APIError {
	return NewAPIError(http.StatusInternalServerError, message, err)
}

func NewBadGatewayError(message string, err error) *APIError {
	return NewAPIError(http.StatusBadGateway, message, err)
}

func NewNotFoundError(message string, err error) *APIError {
	return NewAPIError(http.StatusNotFound, message, err)
}
