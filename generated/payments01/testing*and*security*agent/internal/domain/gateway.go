// file: internal/domain/gateway.go
package domain

import (
	"context"
)

// PaymentGateway defines the interface for payment gateway operations
type PaymentGateway interface {
	// CreatePaymentIntent creates a new payment intent
	CreatePaymentIntent(ctx context.Context, req CreatePaymentIntentRequest) (*PaymentIntentResponse, error)
	
	// ConfirmPaymentIntent confirms a payment intent
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResponse, error)
	
	// CancelPaymentIntent cancels a payment intent
	CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResponse, error)
	
	// CreateRefund creates a refund
	CreateRefund(ctx context.Context, req CreateRefundRequest) (*RefundResponse, error)
	
	// GetPaymentIntent retrieves a payment intent
	GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResponse, error)
	
	// ValidateWebhookSignature validates a webhook signature
	ValidateWebhookSignature(payload []byte, signature string) error
}

// CreatePaymentIntentRequest represents a request to create a payment intent
type CreatePaymentIntentRequest struct {
	Amount              int64
	Currency            string
	PaymentMethodID     string
	CustomerID          string
	Description         string
	Metadata            map[string]string
	IdempotencyKey      string
	CaptureMethod       string // "automatic" or "manual"
}

// PaymentIntentResponse represents a payment intent response from the gateway
type PaymentIntentResponse struct {
	ID                string
	Status            string
	Amount            int64
	Currency          string
	ClientSecret      string
	PaymentMethodID   string
	FailureCode       string
	FailureMessage    string
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentIntentID string
	Amount          int64
	Reason          string
	Metadata        map[string]string
	IdempotencyKey  string
}

// RefundResponse represents a refund response from the gateway
type RefundResponse struct {
	ID            string
	Status        string
	Amount        int64
	Currency      string
	FailureReason string
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	// PublishPaymentCreated publishes a payment created event
	PublishPaymentCreated(ctx context.Context, payment *Payment) error
	
	// PublishPaymentConfirmed publishes a payment confirmed event
	PublishPaymentConfirmed(ctx context.Context, payment *Payment) error
	
	// PublishPaymentFailed publishes a payment failed event
	PublishPaymentFailed(ctx context.Context, payment *Payment) error
	
	// PublishRefundCreated publishes a refund created event
	PublishRefundCreated(ctx context.Context, refund *Refund) error
	
	// PublishRefundCompleted publishes a refund completed event
	PublishRefundCompleted(ctx context.Context, refund *Refund) error
}