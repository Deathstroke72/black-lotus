// file: internal/events/domain/events.go
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventVersion represents the schema version for event evolution
type EventVersion string

const (
	EventVersionV1 EventVersion = "1.0"
	EventVersionV2 EventVersion = "2.0"
)

// BaseEvent contains common fields for all domain events
type BaseEvent struct {
	EventID       string            `json:"event_id"`
	EventType     string            `json:"event_type"`
	EventVersion  EventVersion      `json:"event_version"`
	AggregateID   string            `json:"aggregate_id"`
	AggregateType string            `json:"aggregate_type"`
	CorrelationID string            `json:"correlation_id"`
	CausationID   string            `json:"causation_id,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// NewBaseEvent creates a new base event with generated ID and timestamp
func NewBaseEvent(eventType, aggregateType, aggregateID, correlationID string) BaseEvent {
	return BaseEvent{
		EventID:       uuid.New().String(),
		EventType:     eventType,
		EventVersion:  EventVersionV1,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC(),
		Metadata:      make(map[string]string),
	}
}

// Event type constants
const (
	// Payment events
	EventTypePaymentInitiated        = "payment.initiated"
	EventTypePaymentConfirmed        = "payment.confirmed"
	EventTypePaymentSucceeded        = "payment.succeeded"
	EventTypePaymentFailed           = "payment.failed"
	EventTypePaymentCanceled         = "payment.canceled"
	EventTypePaymentRequiresAction   = "payment.requires_action"

	// Refund events
	EventTypeRefundInitiated         = "refund.initiated"
	EventTypeRefundSucceeded         = "refund.succeeded"
	EventTypeRefundFailed            = "refund.failed"

	// Transaction events
	EventTypeTransactionRecorded     = "transaction.recorded"

	// Payment method events
	EventTypePaymentMethodAdded      = "payment_method.added"
	EventTypePaymentMethodRemoved    = "payment_method.removed"
	EventTypePaymentMethodUpdated    = "payment_method.updated"
)

// Aggregate types
const (
	AggregateTypePayment       = "payment"
	AggregateTypeRefund        = "refund"
	AggregateTypeTransaction   = "transaction"
	AggregateTypePaymentMethod = "payment_method"
)

// =============================================================================
// Payment Events
// =============================================================================

// PaymentInitiatedEvent is published when a payment is first created
type PaymentInitiatedEvent struct {
	BaseEvent
	Payload PaymentInitiatedPayload `json:"payload"`
}

type PaymentInitiatedPayload struct {
	PaymentID       string    `json:"payment_id"`
	CustomerID      string    `json:"customer_id"`
	OrderID         string    `json:"order_id"`
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	PaymentMethodID string    `json:"payment_method_id,omitempty"`
	IdempotencyKey  string    `json:"idempotency_key"`
	Description     string    `json:"description,omitempty"`
	InitiatedAt     time.Time `json:"initiated_at"`
}

func NewPaymentInitiatedEvent(payload PaymentInitiatedPayload, correlationID string) *PaymentInitiatedEvent {
	return &PaymentInitiatedEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentInitiated,
			AggregateTypePayment,
			payload.PaymentID,
			correlationID,
		),
		Payload: payload,
	}
}

// PaymentConfirmedEvent is published when a payment is confirmed by the customer
type PaymentConfirmedEvent struct {
	BaseEvent
	Payload PaymentConfirmedPayload `json:"payload"`
}

type PaymentConfirmedPayload struct {
	PaymentID          string    `json:"payment_id"`
	CustomerID         string    `json:"customer_id"`
	OrderID            string    `json:"order_id"`
	Amount             int64     `json:"amount"`
	Currency           string    `json:"currency"`
	StripePaymentIntent string   `json:"stripe_payment_intent_id"`
	ConfirmedAt        time.Time `json:"confirmed_at"`
}

func NewPaymentConfirmedEvent(payload PaymentConfirmedPayload, correlationID string) *PaymentConfirmedEvent {
	return &PaymentConfirmedEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentConfirmed,
			AggregateTypePayment,
			payload.PaymentID,
			correlationID,
		),
		Payload: payload,
	}
}

// PaymentSucceededEvent is published when a payment is successfully processed
type PaymentSucceededEvent struct {
	BaseEvent
	Payload PaymentSucceededPayload `json:"payload"`
}

type PaymentSucceededPayload struct {
	PaymentID           string    `json:"payment_id"`
	CustomerID          string    `json:"customer_id"`
	OrderID             string    `json:"order_id"`
	Amount              int64     `json:"amount"`
	Currency            string    `json:"currency"`
	StripePaymentIntent string    `json:"stripe_payment_intent_id"`
	StripeChargeID      string    `json:"stripe_charge_id"`
	PaymentMethodID     string    `json:"payment_method_id"`
	CardBrand           string    `json:"card_brand,omitempty"`
	CardLastFour        string    `json:"card_last_four,omitempty"`
	ReceiptURL          string    `json:"receipt_url,omitempty"`
	SucceededAt         time.Time `json:"succeeded_at"`
}

func NewPaymentSucceededEvent(payload PaymentSucceededPayload, correlationID string) *PaymentSucceededEvent {
	return &PaymentSucceededEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentSucceeded,
			AggregateTypePayment,
			payload.PaymentID,
			correlationID,
		),
		Payload: payload,
	}
}

// PaymentFailedEvent is published when a payment fails
type PaymentFailedEvent struct {
	BaseEvent
	Payload PaymentFailedPayload `json:"payload"`
}

type PaymentFailedPayload struct {
	PaymentID           string    `json:"payment_id"`
	CustomerID          string    `json:"customer_id"`
	OrderID             string    `json:"order_id"`
	Amount              int64     `json:"amount"`
	Currency            string    `json:"currency"`
	StripePaymentIntent string    `json:"stripe_payment_intent_id,omitempty"`
	FailureCode         string    `json:"failure_code"`
	FailureMessage      string    `json:"failure_message"`
	Retryable           bool      `json:"retryable"`
	FailedAt            time.Time `json:"failed_at"`
}

func NewPaymentFailedEvent(payload PaymentFailedPayload, correlationID string) *PaymentFailedEvent {
	return &PaymentFailedEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentFailed,
			AggregateTypePayment,
			payload.PaymentID,
			correlationID,
		),
		Payload: payload,
	}
}

// PaymentCanceledEvent is published when a payment is canceled
type PaymentCanceledEvent struct {
	BaseEvent
	Payload PaymentCanceledPayload `json:"payload"`
}

type PaymentCanceledPayload struct {
	PaymentID    string    `json:"payment_id"`
	CustomerID   string    `json:"customer_id"`
	OrderID      string    `json:"order_id"`
	Amount       int64     `json:"amount"`
	Currency     string    `json:"currency"`
	Reason       string    `json:"reason"`
	CanceledBy   string    `json:"canceled_by"` // "customer", "system", "admin"
	CanceledAt   time.Time `json:"canceled_at"`
}

func NewPaymentCanceledEvent(payload PaymentCanceledPayload, correlationID string) *PaymentCanceledEvent {
	return &PaymentCanceledEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentCanceled,
			AggregateTypePayment,
			payload.PaymentID,
			correlationID,
		),
		Payload: payload,
	}
}

// PaymentRequiresActionEvent is published when payment needs additional authentication
type PaymentRequiresActionEvent struct {
	BaseEvent
	Payload PaymentRequiresActionPayload `json:"payload"`
}

type PaymentRequiresActionPayload struct {
	PaymentID           string    `json:"payment_id"`
	CustomerID          string    `json:"customer_id"`
	OrderID             string    `json:"order_id"`
	Amount              int64     `json:"amount"`
	Currency            string    `json:"currency"`
	StripePaymentIntent string    `json:"stripe_payment_intent_id"`
	ActionType          string    `json:"action_type"` // "3ds_authentication", "redirect"
	ClientSecret        string    `json:"client_secret"`
	ExpiresAt           time.Time `json:"expires_at"`
}

func NewPaymentRequiresActionEvent(payload PaymentRequiresActionPayload, correlationID string) *PaymentRequiresActionEvent {
	return &PaymentRequiresActionEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentRequiresAction,
			AggregateTypePayment,
			payload.PaymentID,
			correlationID,
		),
		Payload: payload,
	}
}

// =============================================================================
// Refund Events
// =============================================================================

// RefundInitiatedEvent is published when a refund is requested
type RefundInitiatedEvent struct {
	BaseEvent
	Payload RefundInitiatedPayload `json:"payload"`
}

type RefundInitiatedPayload struct {
	RefundID       string    `json:"refund_id"`
	PaymentID      string    `json:"payment_id"`
	CustomerID     string    `json:"customer_id"`
	OrderID        string    `json:"order_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	Reason         string    `json:"reason"`
	IsPartial      bool      `json:"is_partial"`
	IdempotencyKey string    `json:"idempotency_key"`
	InitiatedBy    string    `json:"initiated_by"` // "customer", "admin", "system"
	InitiatedAt    time.Time `json:"initiated_at"`
}

func NewRefundInitiatedEvent(payload RefundInitiatedPayload, correlationID string) *RefundInitiatedEvent {
	return &RefundInitiatedEvent{
		BaseEvent: NewBaseEvent(
			EventTypeRefundInitiated,
			AggregateTypeRefund,
			payload.RefundID,
			correlationID,
		),
		Payload: payload,
	}
}

// RefundSucceededEvent is published when a refund is processed successfully
type RefundSucceededEvent struct {
	BaseEvent
	Payload RefundSucceededPayload `json:"payload"`
}

type RefundSucceededPayload struct {
	RefundID       string    `json:"refund_id"`
	PaymentID      string    `json:"payment_id"`
	CustomerID     string    `json:"customer_id"`
	OrderID        string    `json:"order_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	StripeRefundID string    `json:"stripe_refund_id"`
	SucceededAt    time.Time `json:"succeeded_at"`
}

func NewRefundSucceededEvent(payload RefundSucceededPayload, correlationID string) *RefundSucceededEvent {
	return &RefundSucceededEvent{
		BaseEvent: NewBaseEvent(
			EventTypeRefundSucceeded,
			AggregateTypeRefund,
			payload.RefundID,
			correlationID,
		),
		Payload: payload,
	}
}

// RefundFailedEvent is published when a refund fails
type RefundFailedEvent struct {
	BaseEvent
	Payload RefundFailedPayload `json:"payload"`
}

type RefundFailedPayload struct {
	RefundID       string    `json:"refund_id"`
	PaymentID      string    `json:"payment_id"`
	CustomerID     string    `json:"customer_id"`
	OrderID        string    `json:"order_id"`
	Amount         int64     `json:"amount"`
	Currency       string    `json:"currency"`
	FailureCode    string    `json:"failure_code"`
	FailureMessage string    `json:"failure_message"`
	FailedAt       time.Time `json:"failed_at"`
}

func NewRefundFailedEvent(payload RefundFailedPayload, correlationID string) *RefundFailedEvent {
	return &RefundFailedEvent{
		BaseEvent: NewBaseEvent(
			EventTypeRefundFailed,
			AggregateTypeRefund,
			payload.RefundID,
			correlationID,
		),
		Payload: payload,
	}
}

// =============================================================================
// Transaction Events
// =============================================================================

// TransactionRecordedEvent is published when a transaction is recorded
type TransactionRecordedEvent struct {
	BaseEvent
	Payload TransactionRecordedPayload `json:"payload"`
}

type TransactionRecordedPayload struct {
	TransactionID   string    `json:"transaction_id"`
	PaymentID       string    `json:"payment_id"`
	RefundID        string    `json:"refund_id,omitempty"`
	CustomerID      string    `json:"customer_id"`
	OrderID         string    `json:"order_id"`
	Type            string    `json:"type"` // "payment", "refund", "chargeback", etc.
	Amount          int64     `json:"amount"`
	Currency        string    `json:"currency"`
	Status          string    `json:"status"`
	BalanceBefore   int64     `json:"balance_before,omitempty"`
	BalanceAfter    int64     `json:"balance_after,omitempty"`
	RecordedAt      time.Time `json:"recorded_at"`
}

func NewTransactionRecordedEvent(payload TransactionRecordedPayload, correlationID string) *TransactionRecordedEvent {
	return &TransactionRecordedEvent{
		BaseEvent: NewBaseEvent(
			EventTypeTransactionRecorded,
			AggregateTypeTransaction,
			payload.TransactionID,
			correlationID,
		),
		Payload: payload,
	}
}

// =============================================================================
// Payment Method Events
// =============================================================================

// PaymentMethodAddedEvent is published when a new payment method is added
type PaymentMethodAddedEvent struct {
	BaseEvent
	Payload PaymentMethodAddedPayload `json:"payload"`
}

type PaymentMethodAddedPayload struct {
	PaymentMethodID string    `json:"payment_method_id"`
	CustomerID      string    `json:"customer_id"`
	Type            string    `json:"type"` // "card", "bank_account", "digital_wallet"
	CardBrand       string    `json:"card_brand,omitempty"`
	CardLastFour    string    `json:"card_last_four,omitempty"`
	CardExpMonth    int       `json:"card_exp_month,omitempty"`
	CardExpYear     int       `json:"card_exp_year,omitempty"`
	IsDefault       bool      `json:"is_default"`
	AddedAt         time.Time `json:"added_at"`
}

func NewPaymentMethodAddedEvent(payload PaymentMethodAddedPayload, correlationID string) *PaymentMethodAddedEvent {
	return &PaymentMethodAddedEvent{
		BaseEvent: NewBaseEvent(
			EventTypePaymentMethodAdded,
			AggregateTypePaymentMethod,
			payload.PaymentMethodID,
			correlationID,
		),
		Payload: payload,
	}
}

// =============================================================================
// Event Serialization
// =============================================================================

// DomainEvent interface for all domain events
type DomainEvent interface {
	GetEventID() string
	GetEventType() string
	GetAggregateID() string
	GetCorrelationID() string
	GetTimestamp() time.Time
	ToJSON() ([]byte, error)
}

func (e BaseEvent) GetEventID() string       { return e.EventID }
func (e BaseEvent) GetEventType() string     { return e.EventType }
func (e BaseEvent) GetAggregateID() string   { return e.AggregateID }
func (e BaseEvent) GetCorrelationID() string { return e.CorrelationID }
func (e BaseEvent) GetTimestamp() time.Time  { return e.Timestamp }

func (e PaymentInitiatedEvent) ToJSON() ([]byte, error)      { return json.Marshal(e) }
func (e PaymentConfirmedEvent) ToJSON() ([]byte, error)      { return json.Marshal(e) }
func (e PaymentSucceededEvent) ToJSON() ([]byte, error)      { return json.Marshal(e) }
func (e PaymentFailedEvent) ToJSON() ([]byte, error)         { return json.Marshal(e) }
func (e PaymentCanceledEvent) ToJSON() ([]byte, error)       { return json.Marshal(e) }
func (e PaymentRequiresActionEvent) ToJSON() ([]byte, error) { return json.Marshal(e) }
func (e RefundInitiatedEvent) ToJSON() ([]byte, error)       { return json.Marshal(e) }
func (e RefundSucceededEvent) ToJSON() ([]byte, error)       { return json.Marshal(e) }
func (e RefundFailedEvent) ToJSON() ([]byte, error)          { return json.Marshal(e) }
func (e TransactionRecordedEvent) ToJSON() ([]byte, error)   { return json.Marshal(e) }
func (e PaymentMethodAddedEvent) ToJSON() ([]byte, error)    { return json.Marshal(e) }