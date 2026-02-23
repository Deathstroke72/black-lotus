// file: internal/domain/event/stock_replenished_event.go
package event

import (
	"time"
)

// StockReplenishedEvent is published when stock is replenished
type StockReplenishedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	MovementID   string                        `json:"movement_id"`
	WarehouseID  string                        `json:"warehouse_id"`
	SupplierID   string                        `json:"supplier_id,omitempty"`
	ReferenceNum string                        `json:"reference_number,omitempty"`
	Items        []StockReplenishedItemDetail  `json:"items"`
}

// StockReplenishedItemDetail contains details of a replenished item
type StockReplenishedItemDetail struct {
	ProductID           string `json:"product_id"`
	SKU                 string `json:"sku"`
	QuantityReplenished int    `json:"quantity_replenished"`
	NewStockLevel       int    `json:"new_stock_level"`
}

// EventName returns the canonical event name
func (e StockReplenishedEvent) EventName() string {
	return "inventory.stock.replenished"
}

// AggregateID returns the aggregate identifier
func (e StockReplenishedEvent) AggregateID() string {
	return e.MovementID
}