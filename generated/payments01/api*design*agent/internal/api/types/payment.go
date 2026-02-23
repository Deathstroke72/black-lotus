// file: internal/api/types/payment.go
package types

import "time"

// PaymentStatus represents the lifecycle state of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusRequiresConfirmation PaymentStatus = "requires_confirmation"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCanceled  PaymentStatus = "canceled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
)

// InitiatePaymentRequest represents the request to create a new payment
// @Description Request payload for initiating a payment
type InitiatePaymentRequest struct {
	// IdempotencyKey ensures exactly-once payment processing
	// Must be unique per payment attempt, typically a UUID
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`

	// OrderID links this payment to an order in the Order Service
	OrderID string `json:"order_id" validate:"required,min=1,max=64"`

	// CustomerID identifies the customer making the payment
	CustomerID string `json:"customer_id" validate:"required,min=1,max=64"`

	// Amount specifies the payment amount and currency
	Amount Money `json:"amount" validate:"required"`

	// PaymentMethodID references a tokenized payment method from Stripe
	// Never contains raw card data (PCI-DSS compliance)
	PaymentMethodID string `json:"payment_method_id" validate:"required,startswith=pm_"`

	// Description appears on customer's statement
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`

	// Metadata for additional context (order details, etc.)
	Metadata map[string]string `json:"metadata,omitempty" validate:"omitempty,dive,keys,max=40,endkeys,max=500"`

	// CaptureMethod determines when funds are captured
	// "automatic" captures immediately, "manual" requires confirmation
	CaptureMethod string `json:"capture_method,omitempty" validate:"omitempty,oneof=automatic manual"`

	// ReturnURL for redirect-based payment methods
	ReturnURL string `json:"return_url,omitempty" validate:"omitempty,url"`
}

// PaymentResponse represents a payment resource
// @Description Payment resource returned by the API
type PaymentResponse struct {
	ID              string            `json:"id"`
	OrderID         string            `json:"order_id"`
	CustomerID      string            `json:"customer_id"`
	Amount          Money             `json:"amount"`
	Status          PaymentStatus     `json:"status"`
	Description     string            `json:"description,omitempty"`
	PaymentMethodID string            `json:"payment_method_id"`
	
	// StripePaymentIntentID for reference (not sensitive)
	StripePaymentIntentID string `json:"stripe_payment_intent_id,omitempty"`
	
	// RefundedAmount tracks total refunded
	RefundedAmount  Money             `json:"refunded_amount,omitempty"`
	
	// FailureCode and FailureMessage for failed payments
	FailureCode    string            `json:"failure_code,omitempty"`
	FailureMessage string            `json:"failure_message,omitempty"`
	
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	ConfirmedAt    *time.Time        `json:"confirmed_at,omitempty"`
}

// ConfirmPaymentRequest for confirming a payment requiring manual capture
type ConfirmPaymentRequest struct {
	// IdempotencyKey for confirmation request
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`

	// AmountToCapture allows partial capture (optional)
	// If omitted, captures full authorized amount
	AmountToCapture *int64 `json:"amount_to_capture,omitempty" validate:"omitempty,gt=0"`
}

// ConfirmPaymentResponse after successful confirmation
type ConfirmPaymentResponse struct {
	Payment PaymentResponse `json:"payment"`
}