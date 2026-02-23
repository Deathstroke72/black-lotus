// file: internal/api/types/errors.go
package types

// ErrorResponse is the standardized error response format
// @Description Standard error response returned by all endpoints
type ErrorResponse struct {
	// Error contains the main error information
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Message is a human-readable error description
	Message string `json:"message"`
	// Details contains additional error context
	Details map[string]interface{} `json:"details,omitempty"`
	// ValidationErrors contains field-level validation errors
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	// TraceID for request tracing
	TraceID string `json:"trace_id,omitempty"`
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// Standard error codes
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInsufficientStock = "INSUFFICIENT_STOCK"
	ErrCodeReservationExpired = "RESERVATION_EXPIRED"
	ErrCodeInvalidState     = "INVALID_STATE"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
)