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