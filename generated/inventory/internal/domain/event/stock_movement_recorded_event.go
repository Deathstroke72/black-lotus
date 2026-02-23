// file: internal/domain/event/stock_movement_recorded_event.go
package event

import (
	"time"
)

// StockMovementRecordedEvent is published for audit trail when any stock movement occurs
type StockMovementRecordedEvent struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	Timestamp     time.Time `json:"timestamp"`
	Version       string    `json:"version"`

	// Payload
	MovementID      string        `json:"movement_id"`
	ProductID       string        `json:"product_id"`
	SKU             string        `json:"sku"`
	WarehouseID     string        `json:"warehouse_id"`
	MovementType    MovementType  `json:"movement_type"`
	Quantity        int           `json:"quantity"`
	PreviousStock   int           `json:"previous_stock"`
	NewStock        int           `json:"new_stock"`
	ReferenceType   string        `json:"reference_type,omitempty"`
	ReferenceID     string        `json:"reference_id,omitempty"`
	Reason          string        `json:"reason,omitempty"`
	PerformedBy     string        `json:"performed_by,omitempty"`
}

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeReservation   MovementType = "RESERVATION"
	MovementTypeRelease       MovementType = "RELEASE"
	MovementTypeDecrement     MovementType = "DECREMENT"
	MovementTypeReplenishment MovementType = "REPLENISHMENT"
	MovementTypeAdjustment    MovementType = "ADJUSTMENT"
	MovementTypeTransfer      MovementType = "TRANSFER"
)

// EventName returns the canonical event name
func (e StockMovementRecordedEvent) EventName() string {
	return "inventory.stock.movement_recorded"
}

// AggregateID returns the aggregate identifier
func (e StockMovementRecordedEvent) AggregateID() string {
	return e.MovementID
}