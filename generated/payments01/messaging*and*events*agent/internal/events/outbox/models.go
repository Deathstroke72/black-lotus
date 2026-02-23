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