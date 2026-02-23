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