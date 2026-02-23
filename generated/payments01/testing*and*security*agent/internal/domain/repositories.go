// file: internal/domain/repositories.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

// PaymentRepository defines the interface for payment persistence
type PaymentRepository interface {
	// Create creates a new payment
	Create(ctx context.Context, payment *Payment) error
	
	// GetByID retrieves a payment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	
	// GetByIdempotencyKey retrieves a payment by idempotency key
	GetByIdempotencyKey(ctx context.Context, key string) (*Payment, error)
	
	// Update updates an existing payment
	Update(ctx context.Context, payment *Payment) error
	
	// GetByCustomerID retrieves payments for a customer
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*Payment, int64, error)
	
	// GetByOrderID retrieves a payment by order ID
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error)
	
	// LockForUpdate locks a payment row for update
	LockForUpdate(ctx context.Context, id uuid.UUID) (*Payment, error)
}

// RefundRepository defines the interface for refund persistence
type RefundRepository interface {
	// Create creates a new refund
	Create(ctx context.Context, refund *Refund) error
	
	// GetByID retrieves a refund by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Refund, error)
	
	// Update updates an existing refund
	Update(ctx context.Context, refund *Refund) error
	
	// GetByPaymentID retrieves refunds for a payment
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*Refund, error)
	
	// GetTotalRefundedAmount gets the total refunded amount for a payment
	GetTotalRefundedAmount(ctx context.Context, paymentID uuid.UUID) (int64, error)
}

// TransactionRepository defines the interface for transaction persistence
type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, tx *Transaction) error
	
	// GetByID retrieves a transaction by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	
	// List retrieves transactions with filters
	List(ctx context.Context, filter TransactionFilter) ([]*Transaction, int64, error)
	
	// GetByPaymentID retrieves transactions for a payment
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*Transaction, error)
}

// PaymentMethodRepository defines the interface for payment method persistence
type PaymentMethodRepository interface {
	// Create creates a new payment method
	Create(ctx context.Context, pm *PaymentMethod) error
	
	// GetByID retrieves a payment method by ID
	GetByID(ctx context.Context, id uuid.UUID) (*PaymentMethod, error)
	
	// GetByCustomerID retrieves payment methods for a customer
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*PaymentMethod, error)
	
	// Update updates an existing payment method
	Update(ctx context.Context, pm *PaymentMethod) error
	
	// Delete soft deletes a payment method
	Delete(ctx context.Context, id uuid.UUID) error
	
	// SetDefault sets a payment method as the default for a customer
	SetDefault(ctx context.Context, customerID, paymentMethodID uuid.UUID) error
}

// IdempotencyRepository defines the interface for idempotency key management
type IdempotencyRepository interface {
	// Store stores an idempotency key with its response
	Store(ctx context.Context, key string, response []byte, ttlSeconds int) error
	
	// Get retrieves a stored response for an idempotency key
	Get(ctx context.Context, key string) ([]byte, bool, error)
	
	// Exists checks if an idempotency key exists
	Exists(ctx context.Context, key string) (bool, error)
}

// WebhookEventRepository defines the interface for webhook event tracking
type WebhookEventRepository interface {
	// Store stores a processed webhook event ID
	Store(ctx context.Context, eventID string) error
	
	// Exists checks if a webhook event has been processed
	Exists(ctx context.Context, eventID string) (bool, error)
}

// UnitOfWork defines the interface for transaction management
type UnitOfWork interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (context.Context, error)
	
	// Commit commits the current transaction
	Commit(ctx context.Context) error
	
	// Rollback rolls back the current transaction
	Rollback(ctx context.Context) error
	
	// WithTransaction executes a function within a transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}