# Messaging & Events Agent Output

# Inventory Microservice - Messaging & Events Layer

I'll provide a comprehensive implementation of the messaging and eventing layer for the inventory microservice.

## 1. Domain Event Structs (Published Events)

```go
// file: internal/events/domain_events.go
package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventVersion constants for schema versioning
const (
	EventVersionV1 = "1.0"
	EventVersionV2 = "2.0"
)

// BaseEvent contains common fields for all domain events
type BaseEvent struct {
	EventID       string            `json:"event_id"`
	EventType     string            `json:"event_type"`
	EventVersion  string            `json:"event_version"`
	AggregateID   string            `json:"aggregate_id"`
	AggregateType string            `json:"aggregate_type"`
	CorrelationID string            `json:"correlation_id"`
	CausationID   string            `json:"causation_id,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// NewBaseEvent creates a new base event with generated ID and timestamp
func NewBaseEvent(eventType, aggregateID, aggregateType, correlationID string) BaseEvent {
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

// --- Published Events ---

// StockReservedEvent is published when stock is successfully reserved for an order
type StockReservedEvent struct {
	BaseEvent
	Payload StockReservedPayload `json:"payload"`
}

type StockReservedPayload struct {
	ReservationID string              `json:"reservation_id"`
	OrderID       string              `json:"order_id"`
	Items         []ReservedItem      `json:"items"`
	ReservedAt    time.Time           `json:"reserved_at"`
	ExpiresAt     time.Time           `json:"expires_at"`
	TotalQuantity int                 `json:"total_quantity"`
}

type ReservedItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	WarehouseID string `json:"warehouse_id"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
}

// StockReservationFailedEvent is published when stock reservation fails
type StockReservationFailedEvent struct {
	BaseEvent
	Payload StockReservationFailedPayload `json:"payload"`
}

type StockReservationFailedPayload struct {
	OrderID       string                `json:"order_id"`
	Reason        string                `json:"reason"`
	FailedItems   []FailedReservation   `json:"failed_items"`
	FailedAt      time.Time             `json:"failed_at"`
}

type FailedReservation struct {
	ProductID        string `json:"product_id"`
	VariantID        string `json:"variant_id,omitempty"`
	SKU              string `json:"sku"`
	RequestedQty     int    `json:"requested_quantity"`
	AvailableQty     int    `json:"available_quantity"`
	Reason           string `json:"reason"`
}

// StockReleasedEvent is published when reserved stock is released
type StockReleasedEvent struct {
	BaseEvent
	Payload StockReleasedPayload `json:"payload"`
}

type StockReleasedPayload struct {
	ReservationID string         `json:"reservation_id"`
	OrderID       string         `json:"order_id"`
	Items         []ReleasedItem `json:"items"`
	Reason        string         `json:"reason"`
	ReleasedAt    time.Time      `json:"released_at"`
}

type ReleasedItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	WarehouseID string `json:"warehouse_id"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
}

// StockDecrementedEvent is published when stock is decremented after fulfillment
type StockDecrementedEvent struct {
	BaseEvent
	Payload StockDecrementedPayload `json:"payload"`
}

type StockDecrementedPayload struct {
	ReservationID   string            `json:"reservation_id"`
	OrderID         string            `json:"order_id"`
	Items           []DecrementedItem `json:"items"`
	DecrementedAt   time.Time         `json:"decremented_at"`
	MovementType    string            `json:"movement_type"`
}

type DecrementedItem struct {
	ProductID       string `json:"product_id"`
	VariantID       string `json:"variant_id,omitempty"`
	WarehouseID     string `json:"warehouse_id"`
	SKU             string `json:"sku"`
	Quantity        int    `json:"quantity"`
	PreviousQty     int    `json:"previous_quantity"`
	NewQty          int    `json:"new_quantity"`
}

// StockReplenishedEvent is published when stock is replenished
type StockReplenishedEvent struct {
	BaseEvent
	Payload StockReplenishedPayload `json:"payload"`
}

type StockReplenishedPayload struct {
	ReplenishmentID string              `json:"replenishment_id"`
	WarehouseID     string              `json:"warehouse_id"`
	Items           []ReplenishedItem   `json:"items"`
	Source          string              `json:"source"`
	Reference       string              `json:"reference,omitempty"`
	ReplenishedAt   time.Time           `json:"replenished_at"`
}

type ReplenishedItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
	PreviousQty int    `json:"previous_quantity"`
	NewQty      int    `json:"new_quantity"`
}

// LowStockAlertEvent is published when stock falls below threshold
type LowStockAlertEvent struct {
	BaseEvent
	Payload LowStockAlertPayload `json:"payload"`
}

type LowStockAlertPayload struct {
	AlertID         string          `json:"alert_id"`
	ProductID       string          `json:"product_id"`
	VariantID       string          `json:"variant_id,omitempty"`
	WarehouseID     string          `json:"warehouse_id"`
	SKU             string          `json:"sku"`
	CurrentQuantity int             `json:"current_quantity"`
	Threshold       int             `json:"threshold"`
	AlertLevel      string          `json:"alert_level"` // "warning", "critical", "out_of_stock"
	DetectedAt      time.Time       `json:"detected_at"`
}

// StockMovementRecordedEvent is published for audit purposes on any stock movement
type StockMovementRecordedEvent struct {
	BaseEvent
	Payload StockMovementRecordedPayload `json:"payload"`
}

type StockMovementRecordedPayload struct {
	MovementID    string    `json:"movement_id"`
	ProductID     string    `json:"product_id"`
	VariantID     string    `json:"variant_id,omitempty"`
	WarehouseID   string    `json:"warehouse_id"`
	SKU           string    `json:"sku"`
	MovementType  string    `json:"movement_type"`
	Quantity      int       `json:"quantity"`
	PreviousQty   int       `json:"previous_quantity"`
	NewQty        int       `json:"new_quantity"`
	Reference     string    `json:"reference,omitempty"`
	Reason        string    `json:"reason,omitempty"`
	PerformedBy   string    `json:"performed_by,omitempty"`
	RecordedAt    time.Time `json:"recorded_at"`
}

// InventorySnapshotEvent is published periodically for reporting/analytics
type InventorySnapshotEvent struct {
	BaseEvent
	Payload InventorySnapshotPayload `json:"payload"`
}

type InventorySnapshotPayload struct {
	SnapshotID    string              `json:"snapshot_id"`
	WarehouseID   string              `json:"warehouse_id,omitempty"`
	Items         []SnapshotItem      `json:"items"`
	TotalProducts int                 `json:"total_products"`
	TotalQuantity int                 `json:"total_quantity"`
	TotalReserved int                 `json:"total_reserved"`
	GeneratedAt   time.Time           `json:"generated_at"`
}

type SnapshotItem struct {
	ProductID       string `json:"product_id"`
	VariantID       string `json:"variant_id,omitempty"`
	SKU             string `json:"sku"`
	WarehouseID     string `json:"warehouse_id"`
	Quantity        int    `json:"quantity"`
	ReservedQty     int    `json:"reserved_quantity"`
	AvailableQty    int    `json:"available_quantity"`
}

// Serialize converts event to JSON bytes
func (e *StockReservedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockReservationFailedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockReleasedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockDecrementedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockReplenishedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *LowStockAlertEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockMovementRecordedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *InventorySnapshotEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}
```

## 2. Consumed Events (From Order Service)

```go
// file: internal/events/consumed_events.go
package events

import (
	"encoding/json"
	"time"
)

// --- Consumed Events from Order Service ---

// OrderCreatedEvent is consumed when a new order is created
type OrderCreatedEvent struct {
	BaseEvent
	Payload OrderCreatedPayload `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID       string           `json:"order_id"`
	CustomerID    string           `json:"customer_id"`
	Items         []OrderItem      `json:"items"`
	ShippingAddr  Address          `json:"shipping_address"`
	OrderedAt     time.Time        `json:"ordered_at"`
	Priority      string           `json:"priority,omitempty"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	VariantID   string  `json:"variant_id,omitempty"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// OrderCancelledEvent is consumed when an order is cancelled
type OrderCancelledEvent struct {
	BaseEvent
	Payload OrderCancelledPayload `json:"payload"`
}

type OrderCancelledPayload struct {
	OrderID      string    `json:"order_id"`
	CustomerID   string    `json:"customer_id"`
	Reason       string    `json:"reason"`
	CancelledBy  string    `json:"cancelled_by"`
	CancelledAt  time.Time `json:"cancelled_at"`
}

// OrderFulfilledEvent is consumed when an order is fulfilled/shipped
type OrderFulfilledEvent struct {
	BaseEvent
	Payload OrderFulfilledPayload `json:"payload"`
}

type OrderFulfilledPayload struct {
	OrderID       string             `json:"order_id"`
	CustomerID    string             `json:"customer_id"`
	Items         []FulfilledItem    `json:"items"`
	ShipmentID    string             `json:"shipment_id"`
	FulfilledAt   time.Time          `json:"fulfilled_at"`
	WarehouseID   string             `json:"warehouse_id"`
}

type FulfilledItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
}

// OrderModifiedEvent is consumed when an order is modified
type OrderModifiedEvent struct {
	BaseEvent
	Payload OrderModifiedPayload `json:"payload"`
}

type OrderModifiedPayload struct {
	OrderID       string          `json:"order_id"`
	CustomerID    string          `json:"customer_id"`
	OldItems      []OrderItem     `json:"old_items"`
	NewItems      []OrderItem     `json:"new_items"`
	ModifiedBy    string          `json:"modified_by"`
	ModifiedAt    time.Time       `json:"modified_at"`
	Reason        string          `json:"reason,omitempty"`
}

// Deserialize methods for consumed events
func DeserializeOrderCreatedEvent(data []byte) (*OrderCreatedEvent, error) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DeserializeOrderCancelledEvent(data []byte) (*OrderCancelledEvent, error) {
	var event OrderCancelledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DeserializeOrderFulfilledEvent(data []byte) (*OrderFulfilledEvent, error) {
	var event OrderFulfilledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DeserializeOrderModifiedEvent(data []byte) (*OrderModifiedEvent, error) {
	var event OrderModifiedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
```

## 3. Outbox Pattern Implementation

```go
// file: internal/outbox/models.go
package outbox

import (
	"time"
)

// OutboxStatus represents the status of an outbox entry
type OutboxStatus string

const (
	StatusPending   OutboxStatus = "pending"
	StatusPublished OutboxStatus = "published"
	StatusFailed    OutboxStatus = "failed"
)

// OutboxEntry represents a message in the outbox table
type OutboxEntry struct {
	ID            string       `db:"id"`
	AggregateType string       `db:"aggregate_type"`
	AggregateID   string       `db:"aggregate_id"`
	EventType     string       `db:"event_type"`
	Topic         string       `db:"topic"`
	Key           string       `db:"key"`
	Payload       []byte       `db:"payload"`
	Headers       []byte       `db:"headers"`
	Status        OutboxStatus `db:"status"`
	RetryCount    int          `db:"retry_count"`
	MaxRetries    int          `db:"max_retries"`
	LastError     *string      `db:"last_error"`
	CreatedAt     time.Time    `db:"created_at"`
	ProcessedAt   *time.Time   `db:"processed_at"`
	ScheduledAt   time.Time    `db:"scheduled_at"`
}
```

```go
// file: internal/outbox/repository.go
package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
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

// InsertInTx inserts an outbox entry within an existing transaction
func (r *Repository) InsertInTx(ctx context.Context, tx *sqlx.Tx, entry *OutboxEntry) error {
	query := `
		INSERT INTO event_outbox (
			id, aggregate_type, aggregate_id, event_type, topic, key,
			payload, headers, status, retry_count, max_retries, created_at, scheduled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.ScheduledAt.IsZero() {
		entry.ScheduledAt = entry.CreatedAt
	}
	if entry.Status == "" {
		entry.Status = StatusPending
	}
	if entry.MaxRetries == 0 {
		entry.MaxRetries = 5
	}

	_, err := tx.ExecContext(ctx, query,
		entry.ID,
		entry.AggregateType,
		entry.AggregateID,
		entry.EventType,
		entry.Topic,
		entry.Key,
		entry.Payload,
		entry.Headers,
		entry.Status,
		entry.RetryCount,
		entry.MaxRetries,
		entry.CreatedAt,
		entry.ScheduledAt,
	)

	return err
}

// FetchPendingBatch fetches a batch of pending outbox entries for processing
func (r *Repository) FetchPendingBatch(ctx context.Context, batchSize int) ([]OutboxEntry, error) {
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, topic, key,
			   payload, headers, status, retry_count, max_retries, last_error,
			   created_at, processed_at, scheduled_at
		FROM event_outbox
		WHERE status = $1
		  AND scheduled_at <= $2
		  AND retry_count < max_retries
		ORDER BY created_at ASC
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`

	var entries []OutboxEntry
	err := r.db.SelectContext(ctx, &entries, query, StatusPending, time.Now().UTC(), batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending entries: %w", err)
	}

	return entries, nil
}

// MarkAsPublished marks an outbox entry as successfully published
func (r *Repository) MarkAsPublished(ctx context.Context, id string) error {
	query := `
		UPDATE event_outbox
		SET status = $1, processed_at = $2
		WHERE id = $3
	`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, StatusPublished, now, id)
	return err
}

// MarkAsFailed marks an outbox entry as failed with retry scheduling
func (r *Repository) MarkAsFailed(ctx context.Context, id string, errMsg string, retryCount int) error {
	// Exponential backoff: 2^retryCount seconds, max 5 minutes
	backoffSeconds := 1 << retryCount
	if backoffSeconds > 300 {
		backoffSeconds = 300
	}
	nextRetry := time.Now().UTC().Add(time.Duration(backoffSeconds) * time.Second)

	query := `
		UPDATE event_outbox
		SET retry_count = $1, last_error = $2, scheduled_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, retryCount, errMsg, nextRetry, id)
	return err
}

// MarkAsDeadLetter marks an entry as permanently failed
func (r *Repository) MarkAsDeadLetter(ctx context.Context, id string, errMsg string) error {
	query := `
		UPDATE event_outbox
		SET status = $1, last_error = $2, processed_at = $3
		WHERE id = $4
	`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, StatusFailed, errMsg, now, id)
	return err
}

// CleanupOldEntries removes successfully processed entries older than retention period
func (r *Repository) CleanupOldEntries(ctx context.Context, retentionDays int) (int64, error) {
	query := `
		DELETE FROM event_outbox
		WHERE status = $1
		  AND processed_at < $2
	`

	cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays)
	result, err := r.db.ExecContext(ctx, query, StatusPublished, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// CreateOutboxEntry is a helper to create an outbox entry from an event
func CreateOutboxEntry(aggregateType, aggregateID, eventType, topic string, payload interface{}, headers map[string]string) (*OutboxEntry, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var headersBytes []byte
	if headers != nil {
		headersBytes, err = json.Marshal(headers)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal headers: %w", err)
		}
	}

	return &OutboxEntry{
		ID:            uuid.New().String(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Topic:         topic,
		Key:           aggregateID, // Use aggregate ID as partition key
		Payload:       payloadBytes,
		Headers:       headersBytes,
		Status:        StatusPending,
		RetryCount:    0,
		MaxRetries:    5,
		CreatedAt:     time.Now().UTC(),
		ScheduledAt:   time.Now().UTC(),
	}, nil
}
```

## 4. Kafka Producer with Transactional Outbox

```go
// file: internal/kafka/producer.go
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"

	"inventory-service/internal/outbox"
)

// ProducerConfig holds configuration for the Kafka producer
type ProducerConfig struct {
	Brokers           []string
	ClientID          string
	RequiredAcks      sarama.RequiredAcks
	MaxRetries        int
	RetryBackoff      time.Duration
	FlushFrequency    time.Duration
	FlushMessages     int
	EnableIdempotent  bool
	CompressionCodec  sarama.CompressionCodec
}

// DefaultProducerConfig returns sensible defaults
func DefaultProducerConfig(brokers []string) ProducerConfig {
	return ProducerConfig{
		Brokers:          brokers,
		ClientID:         "inventory-service",
		RequiredAcks:     sarama.WaitForAll,
		MaxRetries:       5,
		RetryBackoff:     100 * time.Millisecond,
		FlushFrequency:   500 * time.Millisecond,
		FlushMessages:    100,
		EnableIdempotent: true,
		CompressionCodec: sarama.CompressionSnappy,
	}
}

// Producer wraps Sarama async producer with outbox pattern support
type Producer struct {
	producer     sarama.AsyncProducer
	config       ProducerConfig
	logger       zerolog.Logger
	outboxRepo   *outbox.Repository
	
	wg           sync.WaitGroup
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg ProducerConfig, outboxRepo *outbox.Repository, logger zerolog.Logger) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = cfg.ClientID
	saramaConfig.Producer.RequiredAcks = cfg.RequiredAcks
	saramaConfig.Producer.Retry.Max = cfg.MaxRetries
	saramaConfig.Producer.Retry.Backoff = cfg.RetryBackoff
	saramaConfig.Producer.Flush.Frequency = cfg.FlushFrequency
	saramaConfig.Producer.Flush.Messages = cfg.FlushMessages
	saramaConfig.Producer.Idempotent = cfg.EnableIdempotent
	saramaConfig.Producer.Compression = cfg.CompressionCodec
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Net.MaxOpenRequests = 1 // Required for idempotent producer

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	p := &Producer{
		producer:   producer,
		config:     cfg,
		logger:     logger,
		outboxRepo: outboxRepo,
		shutdownCh: make(chan struct{}),
	}

	// Start success/error handlers
	p.wg.Add(2)
	go p.handleSuccesses()
	go p.handleErrors()

	return p, nil
}

// handleSuccesses processes successful message deliveries
func (p *Producer) handleSuccesses() {
	defer p.wg.Done()

	for {
		select {
		case msg, ok := <-p.producer.Successes():
			if !ok {
				return
			}
			
			// Extract outbox ID from metadata
			if outboxID, ok := msg.Metadata.(string); ok {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := p.outboxRepo.MarkAsPublished(ctx, outboxID); err != nil {
					p.logger.Error().
						Err(err).
						Str("outbox_id", outboxID).
						Msg("Failed to mark outbox entry as published")
				} else {
					p.logger.Debug().
						Str("outbox_id", outboxID).
						Str("topic", msg.Topic).
						Int32("partition", msg.Partition).
						Int64("offset", msg.Offset).
						Msg("Message published successfully")
				}
				cancel()
			}

		case <-p.shutdownCh:
			return
		}
	}
}

// handleErrors processes failed message deliveries
func (p *Producer) handleErrors() {
	defer p.wg.Done()

	for {
		select {
		case err, ok := <-p.producer.Errors():
			if !ok {
				return
			}

			// Extract outbox ID and retry count from metadata
			if outboxID, ok := err.Msg.Metadata.(string); ok {
				p.logger.Error().
					Err(err.Err).
					Str("outbox_id", outboxID).
					Str("topic", err.Msg.Topic).
					Msg("Failed to publish message")

				// The outbox processor will handle retries
			}

		case <-p.shutdownCh:
			return
		}
	}
}

// PublishFromOutbox publishes a message from an outbox entry
func (p *Producer) PublishFromOutbox(entry *outbox.OutboxEntry) error {
	// Parse headers
	var headers []sarama.RecordHeader
	if len(entry.Headers) > 0 {
		var headerMap map[string]string
		if err := json.Unmarshal(entry.Headers, &headerMap); err == nil {
			for k, v := range headerMap {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
		}
	}

	// Add standard headers
	headers = append(headers,
		sarama.RecordHeader{Key: []byte("event_type"), Value: []byte(entry.EventType)},
		sarama.RecordHeader{Key: []byte("aggregate_type"), Value: []byte(entry.AggregateType)},
		sarama.RecordHeader{Key: []byte("aggregate_id"), Value: []byte(entry.AggregateID)},
		sarama.RecordHeader{Key: []byte("outbox_id"), Value: []byte(entry.ID)},
	)

	msg := &sarama.ProducerMessage{
		Topic:    entry.Topic,
		Key:      sarama.StringEncoder(entry.Key),
		Value:    sarama.ByteEncoder(entry.Payload),
		Headers:  headers,
		Metadata: entry.ID, // Store outbox ID for success/error handling
	}

	select {
	case p.producer.Input() <- msg:
		return nil
	case <-p.shutdownCh:
		return fmt.Errorf("producer is shutting down")
	}
}

// Close gracefully shuts down the producer
func (p *Producer) Close() error {
	p.shutdownOnce.Do(func() {
		close(p.shutdownCh)
	})

	// Close the producer
	if err := p.producer.Close(); err != nil {
		p.logger.Error().Err(err).Msg("Error closing producer")
		return err
	}

	// Wait for handlers to finish
	p.wg.Wait()
	p.logger.Info().Msg("Producer shut down gracefully")
	return nil
}
```

```go
// file: internal/kafka/outbox_processor.go
package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"inventory-service/internal/outbox"
)

// OutboxProcessor processes outbox entries and publishes them to Kafka
type OutboxProcessor struct {
	outboxRepo   *outbox.Repository
	producer     *Producer
	logger       zerolog.Logger
	
	batchSize     int
	pollInterval  time.Duration
	
	wg           sync.WaitGroup
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
}

// OutboxProcessorConfig holds configuration for the outbox processor
type OutboxProcessorConfig struct {
	BatchSize    int
	PollInterval time.Duration
}

// DefaultOutboxProcessorConfig returns sensible defaults
func DefaultOutboxProcessorConfig() OutboxProcessorConfig {
	return OutboxProcessorConfig