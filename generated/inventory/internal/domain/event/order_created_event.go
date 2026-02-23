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