// file: internal/domain/event/low_stock_alert_event.go
package event

import (
	"time"
)

// LowStockAlertEvent is published when stock falls below minimum threshold
type LowStockAlertEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	AlertID       string              `json:"alert_id"`
	ProductID     string              `json:"product_id"`
	SKU           string              `json:"sku"`
	ProductName   string              `json:"product_name"`
	WarehouseID   string              `json:"warehouse_id"`
	WarehouseName string              `json:"warehouse_name"`
	CurrentStock  int                 `json:"current_stock"`
	MinimumStock  int                 `json:"minimum_stock"`
	Severity      LowStockSeverity    `json:"severity"`
}

// LowStockSeverity represents the severity of a low stock alert
type LowStockSeverity string

const (
	SeverityWarning  LowStockSeverity = "WARNING"
	SeverityCritical LowStockSeverity = "CRITICAL"
	SeverityOutOfStock LowStockSeverity = "OUT_OF_STOCK"
)

// EventName returns the canonical event name
func (e LowStockAlertEvent) EventName() string {
	return "inventory.stock.low_stock_alert"
}

// AggregateID returns the aggregate identifier
func (e LowStockAlertEvent) AggregateID() string {
	return e.AlertID
}