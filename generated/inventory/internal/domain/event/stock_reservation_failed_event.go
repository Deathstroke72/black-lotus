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