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