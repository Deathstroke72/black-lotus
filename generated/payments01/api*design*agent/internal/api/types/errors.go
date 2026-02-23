// file: internal/api/types/errors.go
package types

import (
	"encoding/json"
	"net/http"
)

// ErrorCode represents machine-readable error codes
type ErrorCode string

const (
	// General errors
	ErrCodeInvalidRequest     ErrorCode = "INVALID_REQUEST"
	ErrCodeValidationFailed   ErrorCode = "VALIDATION_FAILED"
	ErrCodeResourceNotFound   ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
	
	// Idempotency errors
	ErrCodeIdempotencyKeyMissing  ErrorCode = "IDEMPOTENCY_KEY_MISSING"
	ErrCodeIdempotencyKeyConflict ErrorCode = "IDEMPOTENCY_KEY_CONFLICT"
	
	// Payment-specific errors
	ErrCodePaymentFailed          ErrorCode = "PAYMENT_FAILED"
	ErrCodePaymentDeclined        ErrorCode = "PAYMENT_DECLINED"
	ErrCodeInsufficientFunds      ErrorCode = "INSUFFICIENT_FUNDS"
	ErrCodeInvalidPaymentMethod   ErrorCode = "INVALID_PAYMENT_METHOD"
	ErrCodePaymentAlreadyConfirmed ErrorCode = "PAYMENT_ALREADY_CONFIRMED"
	ErrCodePaymentNotConfirmable  ErrorCode = "PAYMENT_NOT_CONFIRMABLE"
	ErrCodePaymentExpired         ErrorCode = "PAYMENT_EXPIRED"
	
	// Refund-specific errors
	ErrCodeRefundFailed           ErrorCode = "REFUND_FAILED"
	ErrCodeRefundExceedsPayment   ErrorCode = "REFUND_EXCEEDS_PAYMENT"
	ErrCodePaymentNotRefundable   ErrorCode = "PAYMENT_NOT_REFUNDABLE"
	ErrCodeRefundAlreadyProcessed ErrorCode = "REFUND_ALREADY_PROCESSED"
	
	// Gateway errors
	ErrCodeGatewayTimeout    ErrorCode = "GATEWAY_TIMEOUT"
	ErrCodeGatewayUnavailable ErrorCode = "GATEWAY_UNAVAILABLE"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error information
type ErrorDetail struct {
	// Code is a machine-readable error code
	Code ErrorCode `json:"code"`
	
	// Message is a human-readable error description
	Message string `json:"message"`
	
	// Details provides additional context about the error
	Details map[string]interface{} `json:"details,omitempty"`
	
	// ValidationErrors lists field-level validation failures
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	
	// RequestID for tracing/debugging
	RequestID string `json:"request_id,omitempty"`
	
	// RetryAfter indicates when to retry (for rate limiting)
	RetryAfter *int `json:"retry_after,omitempty"`
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Common error constructors

// NewErrorResponse creates a standard error response
func NewErrorResponse(code ErrorCode, message string, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}
}

// WithDetails adds details to the error response
func (e *ErrorResponse) WithDetails(details map[string]interface{}) *ErrorResponse {
	e.Error.Details = details
	return e
}

// WithValidationErrors adds validation errors to the response
func (e *ErrorResponse) WithValidationErrors(errors []ValidationError) *ErrorResponse {
	e.Error.ValidationErrors = errors
	return e
}

// WriteJSON writes the error response to the http.ResponseWriter
func (e *ErrorResponse) WriteJSON(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(e)
}

// HTTPStatusForErrorCode maps error codes to HTTP status codes
func HTTPStatusForErrorCode(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest, ErrCodeValidationFailed, 
	     ErrCodeIdempotencyKeyMissing:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeResourceNotFound:
		return http.StatusNotFound
	case ErrCodeIdempotencyKeyConflict:
		return http.StatusConflict
	case ErrCodePaymentFailed, ErrCodePaymentDeclined, ErrCodeInsufficientFunds,
	     ErrCodeInvalidPaymentMethod, ErrCodeRefundFailed, 
	     ErrCodeRefundExceedsPayment, ErrCodePaymentNotRefundable,
	     ErrCodePaymentNotConfirmable, ErrCodePaymentAlreadyConfirmed:
		return http.StatusUnprocessableEntity
	case ErrCodeRateLimitExceeded:
		return http.StatusTooManyRequests
	case ErrCodeGatewayTimeout:
		return http.StatusGatewayTimeout
	case ErrCodeGatewayUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}