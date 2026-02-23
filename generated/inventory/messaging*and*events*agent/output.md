# Messaging & Events Agent Output

# Inventory Microservice - Messaging & Events Layer

I'll implement a comprehensive event-driven layer for the inventory microservice following Clean Architecture principles.

## 1. Domain Events (Published by Inventory Service)

```go
// file: internal/domain/event/stock_reserved_event.go
package event

import (
	"time"
)

// StockReservedEvent is published when stock is successfully reserved for an order
type StockReservedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	ReservationID string                    `json:"reservation_id"`
	OrderID       string                    `json:"order_id"`
	WarehouseID   string                    `json:"warehouse_id"`
	Items         []StockReservedItemDetail `json:"items"`
	ExpiresAt     time.Time                 `json:"expires_at"`
}

// StockReservedItemDetail contains details of a reserved item
type StockReservedItemDetail struct {
	ProductID        string `json:"product_id"`
	SKU              string `json:"sku"`
	QuantityReserved int    `json:"quantity_reserved"`
	UnitPrice        int64  `json:"unit_price_cents"`
}

// EventName returns the canonical event name
func (e StockReservedEvent) EventName() string {
	return "inventory.stock.reserved"
}

// AggregateID returns the aggregate identifier
func (e StockReservedEvent) AggregateID() string {
	return e.ReservationID
}
```

```go
// file: internal/domain/event/stock_reservation_failed_event.go
package event

import (
	"time"
)

// StockReservationFailedEvent is published when stock reservation fails
type StockReservationFailedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	OrderID       string                         `json:"order_id"`
	FailureReason string                         `json:"failure_reason"`
	FailedItems   []StockReservationFailedDetail `json:"failed_items"`
}

// StockReservationFailedDetail contains details of items that failed reservation
type StockReservationFailedDetail struct {
	ProductID         string `json:"product_id"`
	SKU               string `json:"sku"`
	RequestedQuantity int    `json:"requested_quantity"`
	AvailableQuantity int    `json:"available_quantity"`
}

// EventName returns the canonical event name
func (e StockReservationFailedEvent) EventName() string {
	return "inventory.stock.reservation_failed"
}

// AggregateID returns the aggregate identifier
func (e StockReservationFailedEvent) AggregateID() string {
	return e.OrderID
}
```

```go
// file: internal/domain/event/stock_released_event.go
package event

import (
	"time"
)

// StockReleasedEvent is published when reserved stock is released
type StockReleasedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	ReservationID string                   `json:"reservation_id"`
	OrderID       string                   `json:"order_id"`
	WarehouseID   string                   `json:"warehouse_id"`
	ReleaseReason string                   `json:"release_reason"`
	Items         []StockReleasedItemDetail `json:"items"`
}

// StockReleasedItemDetail contains details of a released item
type StockReleasedItemDetail struct {
	ProductID        string `json:"product_id"`
	SKU              string `json:"sku"`
	QuantityReleased int    `json:"quantity_released"`
}

// EventName returns the canonical event name
func (e StockReleasedEvent) EventName() string {
	return "inventory.stock.released"
}

// AggregateID returns the aggregate identifier
func (e StockReleasedEvent) AggregateID() string {
	return e.ReservationID
}
```

```go
// file: internal/domain/event/stock_decremented_event.go
package event

import (
	"time"
)

// StockDecrementedEvent is published when stock is decremented after fulfillment
type StockDecrementedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	MovementID    string                       `json:"movement_id"`
	ReservationID string                       `json:"reservation_id"`
	OrderID       string                       `json:"order_id"`
	WarehouseID   string                       `json:"warehouse_id"`
	Items         []StockDecrementedItemDetail `json:"items"`
}

// StockDecrementedItemDetail contains details of a decremented item
type StockDecrementedItemDetail struct {
	ProductID          string `json:"product_id"`
	SKU                string `json:"sku"`
	QuantityDecremented int    `json:"quantity_decremented"`
	RemainingStock     int    `json:"remaining_stock"`
}

// EventName returns the canonical event name
func (e StockDecrementedEvent) EventName() string {
	return "inventory.stock.decremented"
}

// AggregateID returns the aggregate identifier
func (e StockDecrementedEvent) AggregateID() string {
	return e.MovementID
}
```

```go
// file: internal/domain/event/stock_replenished_event.go
package event

import (
	"time"
)

// StockReplenishedEvent is published when stock is replenished
type StockReplenishedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	MovementID   string                        `json:"movement_id"`
	WarehouseID  string                        `json:"warehouse_id"`
	SupplierID   string                        `json:"supplier_id,omitempty"`
	ReferenceNum string                        `json:"reference_number,omitempty"`
	Items        []StockReplenishedItemDetail  `json:"items"`
}

// StockReplenishedItemDetail contains details of a replenished item
type StockReplenishedItemDetail struct {
	ProductID           string `json:"product_id"`
	SKU                 string `json:"sku"`
	QuantityReplenished int    `json:"quantity_replenished"`
	NewStockLevel       int    `json:"new_stock_level"`
}

// EventName returns the canonical event name
func (e StockReplenishedEvent) EventName() string {
	return "inventory.stock.replenished"
}

// AggregateID returns the aggregate identifier
func (e StockReplenishedEvent) AggregateID() string {
	return e.MovementID
}
```

```go
// file: internal/domain/event/low_stock_alert_event.go
package event

import (
	"time"
)

// LowStockAlertEvent is published when stock falls below minimum threshold
type LowStockAlertEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	AlertID       string              `json:"alert_id"`
	ProductID     string              `json:"product_id"`
	SKU           string              `json:"sku"`
	ProductName   string              `json:"product_name"`
	WarehouseID   string              `json:"warehouse_id"`
	WarehouseName string              `json:"warehouse_name"`
	CurrentStock  int                 `json:"current_stock"`
	MinimumStock  int                 `json:"minimum_stock"`
	Severity      LowStockSeverity    `json:"severity"`
}

// LowStockSeverity represents the severity of a low stock alert
type LowStockSeverity string

const (
	SeverityWarning  LowStockSeverity = "WARNING"
	SeverityCritical LowStockSeverity = "CRITICAL"
	SeverityOutOfStock LowStockSeverity = "OUT_OF_STOCK"
)

// EventName returns the canonical event name
func (e LowStockAlertEvent) EventName() string {
	return "inventory.stock.low_stock_alert"
}

// AggregateID returns the aggregate identifier
func (e LowStockAlertEvent) AggregateID() string {
	return e.AlertID
}
```

```go
// file: internal/domain/event/stock_movement_recorded_event.go
package event

import (
	"time"
)

// StockMovementRecordedEvent is published for audit trail when any stock movement occurs
type StockMovementRecordedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	MovementID      string        `json:"movement_id"`
	ProductID       string        `json:"product_id"`
	SKU             string        `json:"sku"`
	WarehouseID     string        `json:"warehouse_id"`
	MovementType    MovementType  `json:"movement_type"`
	Quantity        int           `json:"quantity"`
	PreviousStock   int           `json:"previous_stock"`
	NewStock        int           `json:"new_stock"`
	ReferenceType   string        `json:"reference_type,omitempty"`
	ReferenceID     string        `json:"reference_id,omitempty"`
	Reason          string        `json:"reason,omitempty"`
	PerformedBy     string        `json:"performed_by,omitempty"`
}

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeReservation   MovementType = "RESERVATION"
	MovementTypeRelease       MovementType = "RELEASE"
	MovementTypeDecrement     MovementType = "DECREMENT"
	MovementTypeReplenishment MovementType = "REPLENISHMENT"
	MovementTypeAdjustment    MovementType = "ADJUSTMENT"
	MovementTypeTransfer      MovementType = "TRANSFER"
)

// EventName returns the canonical event name
func (e StockMovementRecordedEvent) EventName() string {
	return "inventory.stock.movement_recorded"
}

// AggregateID returns the aggregate identifier
func (e StockMovementRecordedEvent) AggregateID() string {
	return e.MovementID
}
```

```go
// file: internal/domain/event/event.go
package event

import (
	"time"
)

// DomainEvent is the interface all domain events must implement
type DomainEvent interface {
	EventName() string
	AggregateID() string
}

// EventMetadata contains common metadata for all events
type EventMetadata struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`
}

// NewEventMetadata creates a new EventMetadata with current timestamp
func NewEventMetadata(eventID, correlationID, version string) EventMetadata {
	return EventMetadata{
		EventID:       eventID,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC(),
		Version:       version,
	}
}
```

## 2. Consumed Events (From Order Service)

```go
// file: internal/domain/event/order_created_event.go
package event

import (
	"time"
)

// OrderCreatedEvent is consumed from Order Service to reserve stock
type OrderCreatedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	OrderID     string                   `json:"order_id"`
	CustomerID  string                   `json:"customer_id"`
	Items       []OrderItemDetail        `json:"items"`
	WarehouseID string                   `json:"warehouse_id,omitempty"`
}

// OrderItemDetail contains details of an order item
type OrderItemDetail struct {
	ProductID string `json:"product_id"`
	SKU       string `json:"sku"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price_cents"`
}

// EventName returns the canonical event name
func (e OrderCreatedEvent) EventName() string {
	return "order.created"
}

// AggregateID returns the aggregate identifier
func (e OrderCreatedEvent) AggregateID() string {
	return e.OrderID
}
```

```go
// file: internal/domain/event/order_cancelled_event.go
package event

import (
	"time"
)

// OrderCancelledEvent is consumed from Order Service to release reserved stock
type OrderCancelledEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	OrderID          string `json:"order_id"`
	CancellationReason string `json:"cancellation_reason"`
}

// EventName returns the canonical event name
func (e OrderCancelledEvent) EventName() string {
	return "order.cancelled"
}

// AggregateID returns the aggregate identifier
func (e OrderCancelledEvent) AggregateID() string {
	return e.OrderID
}
```

```go
// file: internal/domain/event/order_fulfilled_event.go
package event

import (
	"time"
)

// OrderFulfilledEvent is consumed from Order Service to decrement stock permanently
type OrderFulfilledEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	OrderID     string    `json:"order_id"`
	FulfilledAt time.Time `json:"fulfilled_at"`
}

// EventName returns the canonical event name
func (e OrderFulfilledEvent) EventName() string {
	return "order.fulfilled"
}

// AggregateID returns the aggregate identifier
func (e OrderFulfilledEvent) AggregateID() string {
	return e.OrderID
}
```

## 3. Application Ports for Event Publishing

```go
// file: internal/application/port/event_publisher.go
package port

import (
	"context"
)

// OutboxEntry represents an event to be published via the outbox pattern
type OutboxEntry struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	CorrelationID string
	CreatedAt     int64
}

// EventPublisher defines the port for publishing domain events
type EventPublisher interface {
	// PublishToOutbox stores an event in the outbox table within the current transaction
	PublishToOutbox(ctx context.Context, entry OutboxEntry) error
}

// OutboxProcessor defines the port for processing outbox entries
type OutboxProcessor interface {
	// ProcessPendingEvents fetches and publishes pending outbox entries
	ProcessPendingEvents(ctx context.Context, batchSize int) error
	// Start begins the background outbox processing
	Start(ctx context.Context) error
	// Stop gracefully stops the outbox processor
	Stop(ctx context.Context) error
}
```

```go
// file: internal/application/port/event_consumer.go
package port

import (
	"context"
)

// ConsumedEvent represents a raw event consumed from Kafka
type ConsumedEvent struct {
	Topic         string
	Partition     int32
	Offset        int64
	Key           []byte
	Value         []byte
	Headers       map[string]string
	Timestamp     int64
}

// EventConsumer defines the port for consuming events
type EventConsumer interface {
	// Start begins consuming events from configured topics
	Start(ctx context.Context) error
	// Stop gracefully stops the consumer
	Stop(ctx context.Context) error
}

// IdempotencyStore defines the port for tracking processed events
type IdempotencyStore interface {
	// IsProcessed checks if an event has already been processed
	IsProcessed(ctx context.Context, eventID string) (bool, error)
	// MarkProcessed marks an event as processed
	MarkProcessed(ctx context.Context, eventID string, topic string) error
}
```

## 4. Kafka Producer with Outbox Pattern

```go
// file: internal/infrastructure/kafka/producer/config.go
package producer

import (
	"time"
)

// Config holds Kafka producer configuration
type Config struct {
	Brokers           []string
	ClientID          string
	Acks              string // "all", "1", "0"
	MaxRetries        int
	RetryBackoff      time.Duration
	MaxBackoff        time.Duration
	BatchSize         int
	LingerMs          int
	CompressionType   string // "none", "gzip", "snappy", "lz4", "zstd"
	IdempotentEnabled bool
}

// DefaultConfig returns a production-ready default configuration
func DefaultConfig() Config {
	return Config{
		Brokers:           []string{"localhost:9092"},
		ClientID:          "inventory-service",
		Acks:              "all",
		MaxRetries:        5,
		RetryBackoff:      100 * time.Millisecond,
		MaxBackoff:        10 * time.Second,
		BatchSize:         16384,
		LingerMs:          5,
		CompressionType:   "snappy",
		IdempotentEnabled: true,
	}
}

// OutboxConfig holds configuration for outbox processing
type OutboxConfig struct {
	PollInterval    time.Duration
	BatchSize       int
	MaxRetries      int
	RetentionPeriod time.Duration // How long to keep published entries
}

// DefaultOutboxConfig returns default outbox configuration
func DefaultOutboxConfig() OutboxConfig {
	return OutboxConfig{
		PollInterval:    100 * time.Millisecond,
		BatchSize:       100,
		MaxRetries:      5,
		RetentionPeriod: 7 * 24 * time.Hour,
	}
}
```

```go
// file: internal/infrastructure/kafka/producer/outbox_publisher.go
package producer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OutboxPublisher reads from the outbox table and publishes to Kafka
type OutboxPublisher struct {
	db            *pgxpool.Pool
	producer      sarama.SyncProducer
	config        OutboxConfig
	kafkaConfig   Config
	topicMapper   TopicMapper
	logger        *slog.Logger
	
	wg            sync.WaitGroup
	stopCh        chan struct{}
	stopped       bool
	mu            sync.RWMutex
}

// TopicMapper maps event types to Kafka topics
type TopicMapper interface {
	GetTopic(eventType string) string
}

// DefaultTopicMapper implements TopicMapper with a prefix-based strategy
type DefaultTopicMapper struct {
	prefix string
}

// NewDefaultTopicMapper creates a new DefaultTopicMapper
func NewDefaultTopicMapper(prefix string) *DefaultTopicMapper {
	return &DefaultTopicMapper{prefix: prefix}
}

// GetTopic returns the Kafka topic for an event type
func (m *DefaultTopicMapper) GetTopic(eventType string) string {
	return fmt.Sprintf("%s.%s", m.prefix, eventType)
}

// OutboxEntry represents a row in the outbox table
type OutboxEntry struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       []byte
	CorrelationID string
	CreatedAt     time.Time
	PublishedAt   *time.Time
	RetryCount    int
	LastError     *string
}

// NewOutboxPublisher creates a new OutboxPublisher
func NewOutboxPublisher(
	db *pgxpool.Pool,
	kafkaConfig Config,
	outboxConfig OutboxConfig,
	topicMapper TopicMapper,
	logger *slog.Logger,
) (*OutboxPublisher, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = kafkaConfig.ClientID
	saramaConfig.Producer.RequiredAcks = parseAcks(kafkaConfig.Acks)
	saramaConfig.Producer.Retry.Max = kafkaConfig.MaxRetries
	saramaConfig.Producer.Retry.Backoff = kafkaConfig.RetryBackoff
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.Idempotent = kafkaConfig.IdempotentEnabled
	saramaConfig.Producer.Compression = parseCompression(kafkaConfig.CompressionType)
	saramaConfig.Net.MaxOpenRequests = 1 // Required for idempotent producer

	producer, err := sarama.NewSyncProducer(kafkaConfig.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &OutboxPublisher{
		db:          db,
		producer:    producer,
		config:      outboxConfig,
		kafkaConfig: kafkaConfig,
		topicMapper: topicMapper,
		logger:      logger,
		stopCh:      make(chan struct{}),
	}, nil
}

func parseAcks(acks string) sarama.RequiredAcks {
	switch acks {
	case "0":
		return sarama.NoResponse
	case "1":
		return sarama.WaitForLocal
	default:
		return sarama.WaitForAll
	}
}

func parseCompression(compression string) sarama.CompressionCodec {
	switch compression {
	case "gzip":
		return sarama.CompressionGZIP
	case "snappy":
		return sarama.CompressionSnappy
	case "lz4":
		return sarama.CompressionLZ4
	case "zstd":
		return sarama.CompressionZSTD
	default:
		return sarama.CompressionNone
	}
}

// Start begins the background outbox processing loop
func (p *OutboxPublisher) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return errors.New("publisher has been stopped")
	}
	p.mu.Unlock()

	p.wg.Add(1)
	go p.processLoop(ctx)

	p.logger.Info("outbox publisher started",
		slog.Duration("poll_interval", p.config.PollInterval),
		slog.Int("batch_size", p.config.BatchSize),
	)

	return nil
}

// Stop gracefully stops the outbox publisher
func (p *OutboxPublisher) Stop(ctx context.Context) error {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return nil
	}
	p.stopped = true
	close(p.stopCh)
	p.mu.Unlock()

	// Wait for processing to complete with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		p.logger.Info("outbox publisher stopped gracefully")
	case <-ctx.Done():
		p.logger.Warn("outbox publisher stop timed out")
	}

	if err := p.producer.Close(); err != nil {
		return fmt.Errorf("failed to close Kafka producer: %w", err)
	}

	return nil
}

func (p *OutboxPublisher) processLoop(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := p.ProcessPendingEvents(ctx, p.config.BatchSize); err != nil {
				p.logger.Error("failed to process pending events",
					slog.String("error", err.Error()),
				)
			}
		}
	}
}

// ProcessPendingEvents fetches and publishes pending outbox entries
func (p *OutboxPublisher) ProcessPendingEvents(ctx context.Context, batchSize int) error {
	entries, err := p.fetchPendingEntries(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to fetch pending entries: %w", err)
	}

	if len(entries) == 0 {
		return nil
	}

	p.logger.Debug("processing outbox entries",
		slog.Int("count", len(entries)),
	)

	for _, entry := range entries {
		if err := p.publishEntry(ctx, entry); err != nil {
			p.logger.Error("failed to publish outbox entry",
				slog.String("entry_id", entry.ID),
				slog.String("event_type", entry.EventType),
				slog.String("error", err.Error()),
			)
			// Record the error but continue processing other entries
			if updateErr := p.recordPublishError(ctx, entry.ID, err); updateErr != nil {
				p.logger.Error("failed to record publish error",
					slog.String("entry_id", entry.ID),
					slog.String("error", updateErr.Error()),
				)
			}
			continue
		}

		if err := p.markAsPublished(ctx, entry.ID); err != nil {
			p.logger.Error("failed to mark entry as published",
				slog.String("entry_id", entry.ID),
				slog.String("error", err.Error()),
			)
		}
	}

	return nil
}

func (p *OutboxPublisher) fetchPendingEntries(ctx context.Context, limit int) ([]OutboxEntry, error) {
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, 
		       correlation_id, created_at, published_at, retry_count, last_error
		FROM outbox
		WHERE published_at IS NULL 
		  AND retry_count < $1
		ORDER BY created_at ASC
		LIMIT $2
		FOR UPDATE SKIP LOCKED
	`

	rows, err := p.db.Query(ctx, query, p.config.MaxRetries, limit)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var entries []OutboxEntry
	for rows.Next() {
		var entry OutboxEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.AggregateType,
			&entry.AggregateID,
			&entry.EventType,
			&entry.Payload,
			&entry.CorrelationID,
			&entry.CreatedAt,
			&entry.PublishedAt,
			&entry.RetryCount,
			&entry.LastError,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (p *OutboxPublisher) publishEntry(ctx context.Context, entry OutboxEntry) error {
	topic := p.topicMapper.GetTopic(entry.EventType)

	// Build message envelope with metadata
	envelope := map[string]interface{}{
		"event_id":       entry.ID,
		"event_type":     entry.EventType,
		"aggregate_type": entry.AggregateType,
		"aggregate_id":   entry.AggregateID,
		"correlation_id": entry.CorrelationID,
		"timestamp":      entry.CreatedAt.Format(time.RFC3339Nano),
		"payload":        json.RawMessage(entry.Payload),
	}

	value, err := json.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("failed to marshal envelope: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(entry.AggregateID),
		Value: sarama.ByteEncoder(value),
		Headers: []sarama.RecordHeader{
			{Key: []byte("event_id"), Value: []byte(entry.ID)},
			{Key: []byte("event_type"), Value: []byte(entry.EventType)},
			{Key: []byte("correlation_id"), Value: []byte(entry.CorrelationID)},
			{Key: []byte("content_type"), Value: []byte("application/json")},
			{Key: []byte("schema_version"), Value: []byte("1.0")},
		},
		Timestamp: entry.CreatedAt,
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return p.handlePublishError(ctx, entry, err)
	}

	p.logger.Debug("published event to Kafka",
		slog.String("entry_id", entry.ID),
		slog.String("topic", topic),
		slog.Int("partition", int(partition)),
		slog.Int64("offset", offset),
	)

	return nil
}

func (p *OutboxPublisher) handlePublishError(ctx context.Context, entry OutboxEntry, err error) error {
	// Calculate exponential backoff for retry
	backoff := p.calculateBackoff(entry.RetryCount)

	p.logger.Warn("publish failed, will retry",
		slog.String("entry_id", entry.ID),
		slog.Int("retry_count", entry.RetryCount),
		slog.Duration("next_backoff", backoff),
		slog.String("error", err.Error()),
	)

	return err
}

func (p *OutboxPublisher) calculateBackoff(retryCount int) time.Duration {
	backoff := p.kafkaConfig.RetryBackoff * time.Duration(1<<uint(retryCount))
	if backoff > p.kafkaConfig.MaxBackoff {
		backoff = p.kafkaConfig.MaxBackoff
	}
	return backoff
}

func (p *OutboxPublisher) markAsPublished(ctx context.Context, entryID string) error {
	query := `
		UPDATE outbox
		SET published_at = NOW()
		WHERE id = $1
	`

	_, err := p.db.Exec(ctx, query, entryID)
	return err
}

func (p *OutboxPublisher