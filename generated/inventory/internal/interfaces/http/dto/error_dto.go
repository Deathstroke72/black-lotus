// file: internal/interfaces/http/dto/error_dto.go
package dto

import "time"

// ErrorResponse represents the standardized error response envelope.
// @Description Standard error response format for all API errors
type ErrorResponse struct {
	// Error contains the error details
	Error ErrorDetail `json:"error"`
	// RequestID is the unique identifier for request tracing
	RequestID string `json:"request_id,omitempty"`
}

// ErrorDetail contains specific error information.
type ErrorDetail struct {
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Message is a human-readable error description
	Message string `json:"message"`
	// Details contains additional error context
	Details []FieldError `json:"details,omitempty"`
	// Timestamp is when the error occurred
	Timestamp time.Time `json:"timestamp"`
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	// Field is the name of the field that failed validation
	Field string `json:"field"`
	// Message describes why validation failed
	Message string `json:"message"`
}

// Common error codes
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInsufficientStock = "INSUFFICIENT_STOCK"
	ErrCodeInvalidState     = "INVALID_STATE"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
)