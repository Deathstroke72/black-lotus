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