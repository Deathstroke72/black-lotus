// file: internal/domain/errors.go
package domain

import "errors"

var (
	// Payment errors
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentAlreadyExists  = errors.New("payment with this idempotency key already exists")
	ErrPaymentAlreadyConfirmed = errors.New("payment already confirmed")
	ErrPaymentNotConfirmable = errors.New("payment is not in a confirmable state")
	ErrPaymentNotRefundable  = errors.New("payment is not in a refundable state")
	ErrInvalidPaymentStatus  = errors.New("invalid payment status")
	
	// Refund errors
	ErrRefundNotFound        = errors.New("refund not found")
	ErrRefundExceedsPayment  = errors.New("refund amount exceeds available payment amount")
	ErrRefundAlreadyProcessed = errors.New("refund already processed")
	
	// Transaction errors
	ErrTransactionNotFound   = errors.New("transaction not found")
	
	// Payment method errors
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrPaymentMethodInactive = errors.New("payment method is not active")
	
	// Validation errors
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrInvalidCurrency       = errors.New("invalid currency")
	ErrInvalidIdempotencyKey = errors.New("invalid idempotency key")
	ErrMissingCustomerID     = errors.New("customer ID is required")
	ErrMissingOrderID        = errors.New("order ID is required")
	
	// Gateway errors
	ErrGatewayUnavailable    = errors.New("payment gateway unavailable")
	ErrGatewayTimeout        = errors.New("payment gateway timeout")
	ErrGatewayRejected       = errors.New("payment rejected by gateway")
	
	// Webhook errors
	ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
	ErrWebhookEventNotFound    = errors.New("webhook event not found")
	ErrDuplicateWebhookEvent   = errors.New("duplicate webhook event")
	
	// Authorization errors
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)