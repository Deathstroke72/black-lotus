# Backend & Database Agent Output

# Payments Microservice - Complete Backend Implementation

I'll provide a comprehensive implementation of the payments microservice with all requested components.

## 1. PostgreSQL Schema

```sql
-- file: migrations/000001_initial_schema.up.sql

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Custom types for payment statuses
CREATE TYPE payment_status AS ENUM (
    'pending',
    'requires_confirmation',
    'processing',
    'succeeded',
    'failed',
    'canceled',
    'refunded',
    'partially_refunded'
);

CREATE TYPE refund_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'canceled'
);

CREATE TYPE transaction_type AS ENUM (
    'payment',
    'refund',
    'chargeback',
    'fee',
    'adjustment'
);

CREATE TYPE transaction_status AS ENUM (
    'pending',
    'completed',
    'failed',
    'reversed'
);

CREATE TYPE payment_method_type AS ENUM (
    'card',
    'bank_account',
    'digital_wallet'
);

CREATE TYPE card_brand AS ENUM (
    'visa',
    'mastercard',
    'amex',
    'discover',
    'unknown'
);

-- Payment Methods table (PCI-DSS compliant - no raw card data)
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    type payment_method_type NOT NULL,
    
    -- Stripe references (we never store actual card data)
    stripe_payment_method_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255),
    
    -- Card metadata (safe to store per PCI-DSS)
    card_brand card_brand,
    card_last_four VARCHAR(4),
    card_exp_month SMALLINT CHECK (card_exp_month >= 1 AND card_exp_month <= 12),
    card_exp_year SMALLINT CHECK (card_exp_year >= 2020),
    card_fingerprint VARCHAR(255), -- For duplicate detection
    
    -- Bank account metadata
    bank_name VARCHAR(255),
    bank_last_four VARCHAR(4),
    
    -- Digital wallet metadata
    wallet_type VARCHAR(50),
    
    billing_name VARCHAR(255),
    billing_email VARCHAR(255),
    billing_address_line1 VARCHAR(255),
    billing_address_line2 VARCHAR(255),
    billing_city VARCHAR(100),
    billing_state VARCHAR(100),
    billing_postal_code VARCHAR(20),
    billing_country VARCHAR(2),
    
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Indexes defined below
    CONSTRAINT valid_card_data CHECK (
        type != 'card' OR (
            card_brand IS NOT NULL AND 
            card_last_four IS NOT NULL AND 
            card_exp_month IS NOT NULL AND 
            card_exp_year IS NOT NULL
        )
    )
);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Idempotency
    idempotency_key VARCHAR(255) NOT NULL,
    
    -- Business references
    customer_id UUID NOT NULL,
    order_id UUID NOT NULL,
    payment_method_id UUID REFERENCES payment_methods(id),
    
    -- Stripe references
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255),
    
    -- Amount details
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    amount_refunded BIGINT NOT NULL DEFAULT 0 CHECK (amount_refunded >= 0),
    
    -- Status tracking
    status payment_status NOT NULL DEFAULT 'pending',
    failure_code VARCHAR(100),
    failure_message TEXT,
    
    -- Metadata
    description TEXT,
    statement_descriptor VARCHAR(22),
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    
    -- Version for optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT unique_idempotency_key UNIQUE (idempotency_key),
    CONSTRAINT valid_refund_amount CHECK (amount_refunded <= amount)
);

-- Refunds table
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Idempotency
    idempotency_key VARCHAR(255) NOT NULL,
    
    -- References
    payment_id UUID NOT NULL REFERENCES payments(id),
    stripe_refund_id VARCHAR(255) UNIQUE,
    
    -- Amount
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    
    -- Status
    status refund_status NOT NULL DEFAULT 'pending',
    failure_code VARCHAR(100),
    failure_message TEXT,
    
    -- Reason
    reason VARCHAR(50),
    description TEXT,
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    
    -- Version for optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT unique_refund_idempotency UNIQUE (idempotency_key)
);

-- Transactions table (immutable audit log)
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- References
    payment_id UUID REFERENCES payments(id),
    refund_id UUID REFERENCES refunds(id),
    
    -- Transaction details
    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    
    -- Amount (can be negative for refunds)
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    
    -- Net amounts after fees
    net_amount BIGINT,
    fee_amount BIGINT DEFAULT 0,
    
    -- External references
    stripe_transaction_id VARCHAR(255),
    stripe_balance_transaction_id VARCHAR(255),
    
    -- Context
    description TEXT,
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    
    -- Transactions are immutable - no updated_at
    CONSTRAINT valid_transaction_ref CHECK (
        (type = 'payment' AND payment_id IS NOT NULL) OR
        (type = 'refund' AND refund_id IS NOT NULL) OR
        (type IN ('chargeback', 'fee', 'adjustment'))
    )
);

-- Idempotency keys table for request deduplication
CREATE TABLE idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    request_path VARCHAR(255) NOT NULL,
    request_hash VARCHAR(64) NOT NULL, -- SHA-256 of request body
    response_code INTEGER,
    response_body JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 hours',
    locked_until TIMESTAMPTZ
);

-- Outbox table for reliable event publishing
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    next_retry_at TIMESTAMPTZ
);

-- Webhook events table for idempotent webhook processing
CREATE TABLE webhook_events (
    id VARCHAR(255) PRIMARY KEY, -- Stripe event ID
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0
);

-- Indexes for payment_methods
CREATE INDEX idx_payment_methods_customer_id ON payment_methods(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_methods_stripe_customer ON payment_methods(stripe_customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_methods_fingerprint ON payment_methods(card_fingerprint) WHERE card_fingerprint IS NOT NULL;

-- Indexes for payments
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);
CREATE INDEX idx_payments_stripe_intent ON payments(stripe_payment_intent_id) WHERE stripe_payment_intent_id IS NOT NULL;
CREATE INDEX idx_payments_customer_status ON payments(customer_id, status);

-- Indexes for refunds
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC);
CREATE INDEX idx_refunds_stripe_id ON refunds(stripe_refund_id) WHERE stripe_refund_id IS NOT NULL;

-- Indexes for transactions
CREATE INDEX idx_transactions_payment_id ON transactions(payment_id) WHERE payment_id IS NOT NULL;
CREATE INDEX idx_transactions_refund_id ON transactions(refund_id) WHERE refund_id IS NOT NULL;
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_status ON transactions(status);

-- Composite index for transaction history queries
CREATE INDEX idx_transactions_history ON transactions(created_at DESC, type, status);

-- Indexes for idempotency_keys
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);

-- Indexes for outbox_events
CREATE INDEX idx_outbox_unpublished ON outbox_events(created_at) WHERE published_at IS NULL;
CREATE INDEX idx_outbox_retry ON outbox_events(next_retry_at) WHERE published_at IS NULL AND next_retry_at IS NOT NULL;

-- Indexes for webhook_events
CREATE INDEX idx_webhook_unprocessed ON webhook_events(created_at) WHERE processed_at IS NULL;

-- Trigger function for updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers
CREATE TRIGGER update_payment_methods_updated_at BEFORE UPDATE ON payment_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to increment version on update (optimistic locking)
CREATE OR REPLACE FUNCTION increment_version()
RETURNS TRIGGER AS $$
BEGIN
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER increment_payments_version BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION increment_version();

CREATE TRIGGER increment_refunds_version BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION increment_version();
```

```sql
-- file: migrations/000001_initial_schema.down.sql

-- Drop triggers first
DROP TRIGGER IF EXISTS increment_refunds_version ON refunds;
DROP TRIGGER IF EXISTS increment_payments_version ON payments;
DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_payment_methods_updated_at ON payment_methods;

-- Drop functions
DROP FUNCTION IF EXISTS increment_version();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in dependency order
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_methods;

-- Drop custom types
DROP TYPE IF EXISTS payment_method_type;
DROP TYPE IF EXISTS card_brand;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS refund_status;
DROP TYPE IF EXISTS payment_status;
```

## 2. Domain Models

```go
// file: internal/domain/models.go
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the lifecycle state of a payment
type PaymentStatus string

const (
	PaymentStatusPending              PaymentStatus = "pending"
	PaymentStatusRequiresConfirmation PaymentStatus = "requires_confirmation"
	PaymentStatusProcessing           PaymentStatus = "processing"
	PaymentStatusSucceeded            PaymentStatus = "succeeded"
	PaymentStatusFailed               PaymentStatus = "failed"
	PaymentStatusCanceled             PaymentStatus = "canceled"
	PaymentStatusRefunded             PaymentStatus = "refunded"
	PaymentStatusPartiallyRefunded    PaymentStatus = "partially_refunded"
)

func (s PaymentStatus) IsTerminal() bool {
	return s == PaymentStatusSucceeded || s == PaymentStatusFailed ||
		s == PaymentStatusCanceled || s == PaymentStatusRefunded
}

func (s PaymentStatus) CanBeConfirmed() bool {
	return s == PaymentStatusPending || s == PaymentStatusRequiresConfirmation
}

func (s PaymentStatus) CanBeRefunded() bool {
	return s == PaymentStatusSucceeded || s == PaymentStatusPartiallyRefunded
}

// RefundStatus represents the state of a refund
type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "pending"
	RefundStatusProcessing RefundStatus = "processing"
	RefundStatusSucceeded  RefundStatus = "succeeded"
	RefundStatusFailed     RefundStatus = "failed"
	RefundStatusCanceled   RefundStatus = "canceled"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeChargeback TransactionType = "chargeback"
	TransactionTypeFee        TransactionType = "fee"
	TransactionTypeAdjustment TransactionType = "adjustment"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// PaymentMethodType represents the type of payment method
type PaymentMethodType string

const (
	PaymentMethodTypeCard         PaymentMethodType = "card"
	PaymentMethodTypeBankAccount  PaymentMethodType = "bank_account"
	PaymentMethodTypeDigitalWallet PaymentMethodType = "digital_wallet"
)

// CardBrand represents card brands
type CardBrand string

const (
	CardBrandVisa       CardBrand = "visa"
	CardBrandMastercard CardBrand = "mastercard"
	CardBrandAmex       CardBrand = "amex"
	CardBrandDiscover   CardBrand = "discover"
	CardBrandUnknown    CardBrand = "unknown"
)

// Money represents a monetary amount
type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// PaymentMethod represents a customer's payment method (PCI-DSS compliant)
type PaymentMethod struct {
	ID                    uuid.UUID         `json:"id"`
	CustomerID            uuid.UUID         `json:"customer_id"`
	Type                  PaymentMethodType `json:"type"`
	StripePaymentMethodID *string           `json:"stripe_payment_method_id,omitempty"`
	StripeCustomerID      *string           `json:"stripe_customer_id,omitempty"`

	// Card metadata (safe to store)
	CardBrand       *CardBrand `json:"card_brand,omitempty"`
	CardLastFour    *string    `json:"card_last_four,omitempty"`
	CardExpMonth    *int       `json:"card_exp_month,omitempty"`
	CardExpYear     *int       `json:"card_exp_year,omitempty"`
	CardFingerprint *string    `json:"-"` // Internal use only

	// Bank account metadata
	BankName     *string `json:"bank_name,omitempty"`
	BankLastFour *string `json:"bank_last_four,omitempty"`

	// Digital wallet
	WalletType *string `json:"wallet_type,omitempty"`

	// Billing info
	BillingName        *string `json:"billing_name,omitempty"`
	BillingEmail       *string `json:"billing_email,omitempty"`
	BillingAddressLine1 *string `json:"billing_address_line1,omitempty"`
	BillingAddressLine2 *string `json:"billing_address_line2,omitempty"`
	BillingCity        *string `json:"billing_city,omitempty"`
	BillingState       *string `json:"billing_state,omitempty"`
	BillingPostalCode  *string `json:"billing_postal_code,omitempty"`
	BillingCountry     *string `json:"billing_country,omitempty"`

	IsDefault bool       `json:"is_default"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
}

// Payment represents a payment transaction
type Payment struct {
	ID              uuid.UUID     `json:"id"`
	IdempotencyKey  string        `json:"idempotency_key"`
	CustomerID      uuid.UUID     `json:"customer_id"`
	OrderID         uuid.UUID     `json:"order_id"`
	PaymentMethodID *uuid.UUID    `json:"payment_method_id,omitempty"`

	StripePaymentIntentID *string `json:"stripe_payment_intent_id,omitempty"`
	StripeCustomerID      *string `json:"stripe_customer_id,omitempty"`

	Amount         int64         `json:"amount"`
	Currency       string        `json:"currency"`
	AmountRefunded int64         `json:"amount_refunded"`
	Status         PaymentStatus `json:"status"`

	FailureCode    *string `json:"failure_code,omitempty"`
	FailureMessage *string `json:"failure_message,omitempty"`

	Description         *string         `json:"description,omitempty"`
	StatementDescriptor *string         `json:"statement_descriptor,omitempty"`
	Metadata            json.RawMessage `json:"metadata,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ConfirmedAt *time.Time `json:"confirmed_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CanceledAt  *time.Time `json:"canceled_at,omitempty"`

	Version int `json:"-"`
}

// RefundableAmount returns the amount that can still be refunded
func (p *Payment) RefundableAmount() int64 {
	return p.Amount - p.AmountRefunded
}

// Refund represents a refund against a payment
type Refund struct {
	ID             uuid.UUID    `json:"id"`
	IdempotencyKey string       `json:"idempotency_key"`
	PaymentID      uuid.UUID    `json:"payment_id"`
	StripeRefundID *string      `json:"stripe_refund_id,omitempty"`

	Amount   int64        `json:"amount"`
	Currency string       `json:"currency"`
	Status   RefundStatus `json:"status"`

	FailureCode    *string `json:"failure_code,omitempty"`
	FailureMessage *string `json:"failure_message,omitempty"`

	Reason      *string         `json:"reason,omitempty"`
	Description *string         `json:"description,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	Version int `json:"-"`
}

// Transaction represents an immutable transaction record
type Transaction struct {
	ID        uuid.UUID       `json:"id"`
	PaymentID *uuid.UUID      `json:"payment_id,omitempty"`
	RefundID  *uuid.UUID      `json:"refund_id,omitempty"`

	Type   TransactionType   `json:"type"`
	Status TransactionStatus `json:"status"`

	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	NetAmount *int64 `json:"net_amount,omitempty"`
	FeeAmount int64  `json:"fee_amount"`

	StripeTransactionID        *string `json:"stripe_transaction_id,omitempty"`
	StripeBalanceTransactionID *string `json:"stripe_balance_transaction_id,omitempty"`

	Description *string         `json:"description,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`

	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// IdempotencyKey for request deduplication
type IdempotencyKey struct {
	Key          string          `json:"key"`
	RequestPath  string          `json:"request_path"`
	RequestHash  string          `json:"request_hash"`
	ResponseCode *int            `json:"response_code,omitempty"`
	ResponseBody json.RawMessage `json:"response_body,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	ExpiresAt    time.Time       `json:"expires_at"`
	LockedUntil  *time.Time      `json:"locked_until,omitempty"`
}

// OutboxEvent for reliable event publishing
type OutboxEvent struct {
	ID            uuid.UUID       `json:"id"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   uuid.UUID       `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
	PublishedAt   *time.Time      `json:"published_at,omitempty"`
	RetryCount    int             `json:"retry_count"`
	LastError     *string         `json:"last_error,omitempty"`
	NextRetryAt   *time.Time      `json:"next_retry_at,omitempty"`
}

// WebhookEvent for idempotent webhook processing
type WebhookEvent struct {
	ID           string          `json:"id"`
	EventType    string          `json:"event_type"`
	Payload      json.RawMessage `json:"payload"`
	ProcessedAt  *time.Time      `json:"processed_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	ErrorMessage *string         `json:"error_message,omitempty"`
	RetryCount   int             `json:"retry_count"`
}
```

```go
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
```

## 3. Repository Interfaces

```go
// file: internal/repository/interfaces.go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/payments/internal/domain"
)

// PaymentRepository defines the interface for payment data access
type PaymentRepository interface {
	// Create creates a new payment
	Create(ctx context.Context, payment *domain.Payment) error

	// GetByID retrieves a payment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)

	// GetByIdempotencyKey retrieves a payment by idempotency key
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error)

	// GetByStripePaymentIntentID retrieves a payment by Stripe payment intent ID
	GetByStripePaymentIntentID(ctx context.Context, intentID string) (*domain.Payment, error)

	// Update updates a payment with optimistic locking
	Update(ctx context.Context, payment *domain.Payment) error

	// UpdateStatus updates payment status atomically
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.PaymentStatus, version int) error

	// IncrementRefundedAmount atomically increments the refunded amount
	IncrementRefundedAmount(ctx context.Context, id uuid.UUID, amount int64, version int) error

	// ListByCustomer retrieves payments for a customer with pagination
	ListByCustomer(ctx context.Context, customerID uuid.UUID, params ListParams) ([]*domain.Payment, int64, error)

	// ListByOrder retrieves payments for an order
	ListByOrder(ctx context.Context, orderID uuid.UUID) ([]*domain.Payment, error)

	// GetForUpdate retrieves a payment with a row-level lock
	GetForUpdate(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
}

// RefundRepository defines the interface for refund data access
type RefundRepository interface {
	// Create creates a new refund
	Create(ctx context.Context, refund *domain.Refund) error

	// GetByID retrieves a refund by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Refund, error)

	// GetByIdempotencyKey retrieves a refund by idempotency key
	GetByIdempotencyKey(ctx context.Context, key string) (*domain.Refund, error)

	// GetByStripeRefundID retrieves a refund by Stripe refund ID
	GetByStripeRefundID(ctx context.Context, refundID string) (*domain.Refund, error)

	// Update updates a refund with optimistic locking
	Update(ctx context.Context, refund *domain.Refund) error

	// UpdateStatus updates refund status atomically
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RefundStatus, version int) error

	// ListByPayment retrieves refunds for a payment
	ListByPayment(ctx context.Context, paymentID uuid.UUID, params ListParams) ([]*domain.Refund, int64, error)
}

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, txn *domain.Transaction) error

	// GetByID retrieves a transaction by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)

	// UpdateStatus updates transaction status
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.TransactionStatus) error

	// List retrieves transactions with filtering and pagination
	List(ctx context.Context, filter TransactionFilter, params ListParams) ([]*domain.Transaction, int64, error)

	// ListByPayment retrieves transactions for a payment
	ListByPayment(ctx context.Context, paymentID uuid.UUID) ([]*domain.Transaction, error)
}

// PaymentMethodRepository defines the interface for payment method data access
type PaymentMethodRepository interface {
	// Create creates a new payment method
	Create(ctx context.Context, pm *domain.PaymentMethod) error

	// GetByID retrieves a payment method by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PaymentMethod, error)

	// GetByStripeID retrieves a payment method by Stripe payment method ID
	GetByStripeID(ctx context.Context, stripeID string) (*domain.PaymentMethod, error)

	// Update updates a payment method
	Update(ctx context.Context, pm *domain