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