// file: internal/domain/errors.go
package domain

import (
	"errors"
	"fmt"
)

var (
	// General errors
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrInvalidInput       = errors.New("invalid input")
	ErrConcurrencyConflict = errors.New("concurrent modification detected")

	// Payment errors
	ErrPaymentNotFound        = fmt.Errorf("payment %w", ErrNotFound)
	ErrPaymentAlreadyExists   = fmt.Errorf("payment %w", ErrAlreadyExists)
	ErrPaymentInvalidStatus   = errors.New("payment has invalid status for this operation")
	ErrPaymentCannotConfirm   = errors.New("payment cannot be confirmed in current state")
	ErrPaymentCannotRefund    = errors.New("payment cannot be refunded in current state")
	ErrPaymentAlreadyRefunded = errors.New("payment has already been fully refunded")

	// Refund errors
	ErrRefundNotFound       = fmt.Errorf("refund %w", ErrNotFound)
	ErrRefundExceedsAmount  = errors.New("refund amount exceeds refundable amount")
	ErrRefundAlreadyExists  = fmt.Errorf("refund %w", ErrAlreadyExists)

	// Payment method errors
	ErrPaymentMethodNotFound = fmt.Errorf("payment method %w", ErrNotFound)
	ErrPaymentMethodInactive = errors.New("payment method is inactive")

	// Idempotency errors
	ErrIdempotencyKeyLocked   = errors.New("idempotency key is locked by another request")
	ErrIdempotencyKeyMismatch = errors.New("request body does not match original request")

	// External service errors
	ErrStripeError        = errors.New("stripe API error")
	ErrStripeTimeout      = errors.New("stripe API timeout")
	ErrStripeRateLimited  = errors.New("stripe API rate limited")
	ErrWebhookInvalid     = errors.New("invalid webhook signature")
	ErrWebhookDuplicate   = errors.New("webhook event already processed")
)

// DomainError wraps domain errors with additional context
type DomainError struct {
	Err     error
	Message string
	Code    string
	Details map[string]interface{}
}

func (e *DomainError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Err.Error()
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(err error, message, code string) *DomainError {
	return &DomainError{
		Err:     err,
		Message: message,
		Code:    code,
		Details: make(map[string]interface{}),
	}
}

// WithDetail adds a detail to the error
func (e *DomainError) WithDetail(key string, value interface{}) *DomainError {
	e.Details[key] = value
	return e
}