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