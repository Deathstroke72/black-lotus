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