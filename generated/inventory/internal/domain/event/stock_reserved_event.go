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