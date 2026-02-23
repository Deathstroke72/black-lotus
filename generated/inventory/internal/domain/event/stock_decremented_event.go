// file: internal/domain/event/stock_decremented_event.go
package event

import (
	"time"
)

// StockDecrementedEvent is published when stock is decremented after fulfillment
type StockDecrementedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	MovementID    string                       `json:"movement_id"`
	ReservationID string                       `json:"reservation_id"`
	OrderID       string                       `json:"order_id"`
	WarehouseID   string                       `json:"warehouse_id"`
	Items         []StockDecrementedItemDetail `json:"items"`
}

// StockDecrementedItemDetail contains details of a decremented item
type StockDecrementedItemDetail struct {
	ProductID          string `json:"product_id"`
	SKU                string `json:"sku"`
	QuantityDecremented int    `json:"quantity_decremented"`
	RemainingStock     int    `json:"remaining_stock"`
}

// EventName returns the canonical event name
func (e StockDecrementedEvent) EventName() string {
	return "inventory.stock.decremented"
}

// AggregateID returns the aggregate identifier
func (e StockDecrementedEvent) AggregateID() string {
	return e.MovementID
}