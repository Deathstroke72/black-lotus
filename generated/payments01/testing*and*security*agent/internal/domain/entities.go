// file: internal/domain/entities.go
package domain

import (
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

// RefundStatus represents the lifecycle state of a refund
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

// Money represents a monetary amount
type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// Payment represents a payment entity
type Payment struct {
	ID                    uuid.UUID     `json:"id"`
	IdempotencyKey        string        `json:"idempotency_key"`
	CustomerID            uuid.UUID     `json:"customer_id"`
	OrderID               uuid.UUID     `json:"order_id"`
	PaymentMethodID       uuid.UUID     `json:"payment_method_id"`
	Amount                Money         `json:"amount"`
	Status                PaymentStatus `json:"status"`
	StripePaymentIntentID string        `json:"stripe_payment_intent_id,omitempty"`
	Description           string        `json:"description,omitempty"`
	Metadata              map[string]string `json:"metadata,omitempty"`
	FailureCode           string        `json:"failure_code,omitempty"`
	FailureMessage        string        `json:"failure_message,omitempty"`
	RefundedAmount        int64         `json:"refunded_amount"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	ConfirmedAt           *time.Time    `json:"confirmed_at,omitempty"`
}

// Refund represents a refund entity
type Refund struct {
	ID             uuid.UUID    `json:"id"`
	PaymentID      uuid.UUID    `json:"payment_id"`
	Amount         Money        `json:"amount"`
	Status         RefundStatus `json:"status"`
	Reason         string       `json:"reason,omitempty"`
	StripeRefundID string       `json:"stripe_refund_id,omitempty"`
	FailureReason  string       `json:"failure_reason,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	ProcessedAt    *time.Time   `json:"processed_at,omitempty"`
}

// Transaction represents a transaction record
type Transaction struct {
	ID              uuid.UUID         `json:"id"`
	PaymentID       uuid.UUID         `json:"payment_id"`
	RefundID        *uuid.UUID        `json:"refund_id,omitempty"`
	Type            TransactionType   `json:"type"`
	Status          TransactionStatus `json:"status"`
	Amount          Money             `json:"amount"`
	BalanceBefore   int64             `json:"balance_before"`
	BalanceAfter    int64             `json:"balance_after"`
	Description     string            `json:"description,omitempty"`
	StripeChargeID  string            `json:"stripe_charge_id,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

// PaymentMethod represents a saved payment method
type PaymentMethod struct {
	ID                    uuid.UUID         `json:"id"`
	CustomerID            uuid.UUID         `json:"customer_id"`
	Type                  PaymentMethodType `json:"type"`
	StripePaymentMethodID string            `json:"stripe_payment_method_id"`
	CardBrand             string            `json:"card_brand,omitempty"`
	CardLastFour          string            `json:"card_last_four,omitempty"`
	CardExpMonth          int               `json:"card_exp_month,omitempty"`
	CardExpYear           int               `json:"card_exp_year,omitempty"`
	IsDefault             bool              `json:"is_default"`
	IsActive              bool              `json:"is_active"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

// TransactionFilter represents filters for transaction queries
type TransactionFilter struct {
	CustomerID  *uuid.UUID
	PaymentID   *uuid.UUID
	Type        *TransactionType
	Status      *TransactionStatus
	StartDate   *time.Time
	EndDate     *time.Time
	MinAmount   *int64
	MaxAmount   *int64
	Currency    *string
	Limit       int
	Offset      int
}

// WebhookEvent represents a Stripe webhook event
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}