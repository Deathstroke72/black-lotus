# Messaging & Events Agent Output

# Payments Microservice - Complete Messaging & Events Implementation

I'll provide a comprehensive implementation of the messaging/eventing layer for the payments microservice.

## 1. Domain Event Structs (Published Events)

```go
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
```

## 2. Consumed Events (from Order Service)

```go
// file: internal/events/consumed/order_events.go
package consumed

import (
	"encoding/json"
	"time"
)

// Order Service Events that Payments service consumes

const (
	EventTypeOrderCreated   = "order.created"
	EventTypeOrderConfirmed = "order.confirmed"
	EventTypeOrderCanceled  = "order.canceled"
	EventTypeOrderCompleted = "order.completed"
)

// OrderCreatedEvent is consumed to initiate payment processing
type OrderCreatedEvent struct {
	EventID       string              `json:"event_id"`
	EventType     string              `json:"event_type"`
	EventVersion  string              `json:"event_version"`
	AggregateID   string              `json:"aggregate_id"`
	CorrelationID string              `json:"correlation_id"`
	Timestamp     time.Time           `json:"timestamp"`
	Payload       OrderCreatedPayload `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID          string           `json:"order_id"`
	CustomerID       string           `json:"customer_id"`
	TotalAmount      int64            `json:"total_amount"`
	Currency         string           `json:"currency"`
	Items            []OrderItem      `json:"items"`
	ShippingAddress  Address          `json:"shipping_address"`
	BillingAddress   Address          `json:"billing_address"`
	PaymentMethodID  string           `json:"payment_method_id,omitempty"`
	IdempotencyKey   string           `json:"idempotency_key"`
	RequireCapture   bool             `json:"require_capture"` // true = authorize only, false = capture immediately
	Metadata         map[string]string `json:"metadata,omitempty"`
}

type OrderItem struct {
	ItemID      string `json:"item_id"`
	ProductID   string `json:"product_id"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
	TotalPrice  int64  `json:"total_price"`
}

type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// OrderCanceledEvent is consumed to cancel/refund payment
type OrderCanceledEvent struct {
	EventID       string               `json:"event_id"`
	EventType     string               `json:"event_type"`
	EventVersion  string               `json:"event_version"`
	AggregateID   string               `json:"aggregate_id"`
	CorrelationID string               `json:"correlation_id"`
	Timestamp     time.Time            `json:"timestamp"`
	Payload       OrderCanceledPayload `json:"payload"`
}

type OrderCanceledPayload struct {
	OrderID       string    `json:"order_id"`
	CustomerID    string    `json:"customer_id"`
	Reason        string    `json:"reason"`
	CanceledBy    string    `json:"canceled_by"`
	RefundAmount  int64     `json:"refund_amount"` // 0 if no refund needed
	CanceledAt    time.Time `json:"canceled_at"`
}

// OrderConfirmedEvent is consumed to capture authorized payment
type OrderConfirmedEvent struct {
	EventID       string                `json:"event_id"`
	EventType     string                `json:"event_type"`
	EventVersion  string                `json:"event_version"`
	AggregateID   string                `json:"aggregate_id"`
	CorrelationID string                `json:"correlation_id"`
	Timestamp     time.Time             `json:"timestamp"`
	Payload       OrderConfirmedPayload `json:"payload"`
}

type OrderConfirmedPayload struct {
	OrderID      string    `json:"order_id"`
	CustomerID   string    `json:"customer_id"`
	ConfirmedAt  time.Time `json:"confirmed_at"`
}

// ParseOrderCreatedEvent parses raw JSON into OrderCreatedEvent
func ParseOrderCreatedEvent(data []byte) (*OrderCreatedEvent, error) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// ParseOrderCanceledEvent parses raw JSON into OrderCanceledEvent
func ParseOrderCanceledEvent(data []byte) (*OrderCanceledEvent, error) {
	var event OrderCanceledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// ParseOrderConfirmedEvent parses raw JSON into OrderConfirmedEvent
func ParseOrderConfirmedEvent(data []byte) (*OrderConfirmedEvent, error) {
	var event OrderConfirmedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
```

## 3. Outbox Repository and Models

```go
// file: internal/events/outbox/models.go
package outbox

import (
	"time"
)

// OutboxStatus represents the status of an outbox entry
type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusSent       OutboxStatus = "sent"
	OutboxStatusFailed     OutboxStatus = "failed"
)

// OutboxEntry represents a row in the event_outbox table
type OutboxEntry struct {
	ID              string       `db:"id"`
	AggregateType   string       `db:"aggregate_type"`
	AggregateID     string       `db:"aggregate_id"`
	EventType       string       `db:"event_type"`
	EventVersion    string       `db:"event_version"`
	Payload         []byte       `db:"payload"`
	CorrelationID   string       `db:"correlation_id"`
	CausationID     *string      `db:"causation_id"`
	Topic           string       `db:"topic"`
	PartitionKey    string       `db:"partition_key"`
	Status          OutboxStatus `db:"status"`
	RetryCount      int          `db:"retry_count"`
	MaxRetries      int          `db:"max_retries"`
	LastError       *string      `db:"last_error"`
	ScheduledAt     time.Time    `db:"scheduled_at"`
	ProcessedAt     *time.Time   `db:"processed_at"`
	CreatedAt       time.Time    `db:"created_at"`
}

// NewOutboxEntry creates a new outbox entry
func NewOutboxEntry(
	aggregateType, aggregateID, eventType, eventVersion string,
	payload []byte,
	correlationID, topic, partitionKey string,
) *OutboxEntry {
	return &OutboxEntry{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		EventVersion:  string(eventVersion),
		Payload:       payload,
		CorrelationID: correlationID,
		Topic:         topic,
		PartitionKey:  partitionKey,
		Status:        OutboxStatusPending,
		RetryCount:    0,
		MaxRetries:    5,
		ScheduledAt:   time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	}
}
```

```go
// file: internal/events/outbox/repository.go
package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Repository handles outbox persistence operations
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new outbox repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// InsertWithTx inserts an outbox entry within an existing transaction
func (r *Repository) InsertWithTx(ctx context.Context, tx *sqlx.Tx, entry *OutboxEntry) error {
	query := `
		INSERT INTO event_outbox (
			aggregate_type, aggregate_id, event_type, event_version,
			payload, correlation_id, causation_id, topic, partition_key,
			status, retry_count, max_retries, scheduled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at`

	return tx.QueryRowxContext(ctx, query,
		entry.AggregateType,
		entry.AggregateID,
		entry.EventType,
		entry.EventVersion,
		entry.Payload,
		entry.CorrelationID,
		entry.CausationID,
		entry.Topic,
		entry.PartitionKey,
		entry.Status,
		entry.RetryCount,
		entry.MaxRetries,
		entry.ScheduledAt,
	).Scan(&entry.ID, &entry.CreatedAt)
}

// FetchPendingBatch fetches a batch of pending outbox entries for processing
// Uses SELECT FOR UPDATE SKIP LOCKED to allow concurrent processors
func (r *Repository) FetchPendingBatch(ctx context.Context, batchSize int) ([]*OutboxEntry, error) {
	query := `
		SELECT 
			id, aggregate_type, aggregate_id, event_type, event_version,
			payload, correlation_id, causation_id, topic, partition_key,
			status, retry_count, max_retries, last_error, scheduled_at,
			processed_at, created_at
		FROM event_outbox
		WHERE status IN ('pending', 'failed')
		  AND scheduled_at <= NOW()
		  AND retry_count < max_retries
		ORDER BY scheduled_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	var entries []*OutboxEntry
	err := r.db.SelectContext(ctx, &entries, query, batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending outbox entries: %w", err)
	}
	return entries, nil
}

// MarkAsProcessing marks entries as being processed
func (r *Repository) MarkAsProcessing(ctx context.Context, tx *sqlx.Tx, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
		UPDATE event_outbox 
		SET status = 'processing'
		WHERE id IN (?)`, ids)
	if err != nil {
		return fmt.Errorf("failed to build IN query: %w", err)
	}

	query = tx.Rebind(query)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to mark entries as processing: %w", err)
	}
	return nil
}

// MarkAsSent marks an entry as successfully sent
func (r *Repository) MarkAsSent(ctx context.Context, id string) error {
	now := time.Now().UTC()
	query := `
		UPDATE event_outbox 
		SET status = 'sent', processed_at = $2
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return fmt.Errorf("failed to mark entry as sent: %w", err)
	}
	return nil
}

// MarkAsFailed marks an entry as failed and schedules retry with backoff
func (r *Repository) MarkAsFailed(ctx context.Context, id string, errorMsg string) error {
	// Exponential backoff: 2^retry_count seconds, max 5 minutes
	query := `
		UPDATE event_outbox 
		SET 
			status = CASE 
				WHEN retry_count + 1 >= max_retries THEN 'failed'::outbox_status 
				ELSE 'pending'::outbox_status 
			END,
			retry_count = retry_count + 1,
			last_error = $2,
			scheduled_at = NOW() + (LEAST(POWER(2, retry_count + 1), 300) || ' seconds')::interval
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, errorMsg)
	if err != nil {
		return fmt.Errorf("failed to mark entry as failed: %w", err)
	}
	return nil
}

// CleanupSentEntries removes old sent entries (data retention)
func (r *Repository) CleanupSentEntries(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-olderThan)
	query := `DELETE FROM event_outbox WHERE status = 'sent' AND processed_at < $1`

	result, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup sent entries: %