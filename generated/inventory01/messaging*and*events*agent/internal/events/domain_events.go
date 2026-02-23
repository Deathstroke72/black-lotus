// file: internal/events/domain_events.go
package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventVersion constants for schema versioning
const (
	EventVersionV1 = "1.0"
	EventVersionV2 = "2.0"
)

// BaseEvent contains common fields for all domain events
type BaseEvent struct {
	EventID       string            `json:"event_id"`
	EventType     string            `json:"event_type"`
	EventVersion  string            `json:"event_version"`
	AggregateID   string            `json:"aggregate_id"`
	AggregateType string            `json:"aggregate_type"`
	CorrelationID string            `json:"correlation_id"`
	CausationID   string            `json:"causation_id,omitempty"`
	Timestamp     time.Time         `json:"timestamp"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

// NewBaseEvent creates a new base event with generated ID and timestamp
func NewBaseEvent(eventType, aggregateID, aggregateType, correlationID string) BaseEvent {
	return BaseEvent{
		EventID:       uuid.New().String(),
		EventType:     eventType,
		EventVersion:  EventVersionV1,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC(),
		Metadata:      make(map[string]string),
	}
}

// --- Published Events ---

// StockReservedEvent is published when stock is successfully reserved for an order
type StockReservedEvent struct {
	BaseEvent
	Payload StockReservedPayload `json:"payload"`
}

type StockReservedPayload struct {
	ReservationID string              `json:"reservation_id"`
	OrderID       string              `json:"order_id"`
	Items         []ReservedItem      `json:"items"`
	ReservedAt    time.Time           `json:"reserved_at"`
	ExpiresAt     time.Time           `json:"expires_at"`
	TotalQuantity int                 `json:"total_quantity"`
}

type ReservedItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	WarehouseID string `json:"warehouse_id"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
}

// StockReservationFailedEvent is published when stock reservation fails
type StockReservationFailedEvent struct {
	BaseEvent
	Payload StockReservationFailedPayload `json:"payload"`
}

type StockReservationFailedPayload struct {
	OrderID       string                `json:"order_id"`
	Reason        string                `json:"reason"`
	FailedItems   []FailedReservation   `json:"failed_items"`
	FailedAt      time.Time             `json:"failed_at"`
}

type FailedReservation struct {
	ProductID        string `json:"product_id"`
	VariantID        string `json:"variant_id,omitempty"`
	SKU              string `json:"sku"`
	RequestedQty     int    `json:"requested_quantity"`
	AvailableQty     int    `json:"available_quantity"`
	Reason           string `json:"reason"`
}

// StockReleasedEvent is published when reserved stock is released
type StockReleasedEvent struct {
	BaseEvent
	Payload StockReleasedPayload `json:"payload"`
}

type StockReleasedPayload struct {
	ReservationID string         `json:"reservation_id"`
	OrderID       string         `json:"order_id"`
	Items         []ReleasedItem `json:"items"`
	Reason        string         `json:"reason"`
	ReleasedAt    time.Time      `json:"released_at"`
}

type ReleasedItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	WarehouseID string `json:"warehouse_id"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
}

// StockDecrementedEvent is published when stock is decremented after fulfillment
type StockDecrementedEvent struct {
	BaseEvent
	Payload StockDecrementedPayload `json:"payload"`
}

type StockDecrementedPayload struct {
	ReservationID   string            `json:"reservation_id"`
	OrderID         string            `json:"order_id"`
	Items           []DecrementedItem `json:"items"`
	DecrementedAt   time.Time         `json:"decremented_at"`
	MovementType    string            `json:"movement_type"`
}

type DecrementedItem struct {
	ProductID       string `json:"product_id"`
	VariantID       string `json:"variant_id,omitempty"`
	WarehouseID     string `json:"warehouse_id"`
	SKU             string `json:"sku"`
	Quantity        int    `json:"quantity"`
	PreviousQty     int    `json:"previous_quantity"`
	NewQty          int    `json:"new_quantity"`
}

// StockReplenishedEvent is published when stock is replenished
type StockReplenishedEvent struct {
	BaseEvent
	Payload StockReplenishedPayload `json:"payload"`
}

type StockReplenishedPayload struct {
	ReplenishmentID string              `json:"replenishment_id"`
	WarehouseID     string              `json:"warehouse_id"`
	Items           []ReplenishedItem   `json:"items"`
	Source          string              `json:"source"`
	Reference       string              `json:"reference,omitempty"`
	ReplenishedAt   time.Time           `json:"replenished_at"`
}

type ReplenishedItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
	PreviousQty int    `json:"previous_quantity"`
	NewQty      int    `json:"new_quantity"`
}

// LowStockAlertEvent is published when stock falls below threshold
type LowStockAlertEvent struct {
	BaseEvent
	Payload LowStockAlertPayload `json:"payload"`
}

type LowStockAlertPayload struct {
	AlertID         string          `json:"alert_id"`
	ProductID       string          `json:"product_id"`
	VariantID       string          `json:"variant_id,omitempty"`
	WarehouseID     string          `json:"warehouse_id"`
	SKU             string          `json:"sku"`
	CurrentQuantity int             `json:"current_quantity"`
	Threshold       int             `json:"threshold"`
	AlertLevel      string          `json:"alert_level"` // "warning", "critical", "out_of_stock"
	DetectedAt      time.Time       `json:"detected_at"`
}

// StockMovementRecordedEvent is published for audit purposes on any stock movement
type StockMovementRecordedEvent struct {
	BaseEvent
	Payload StockMovementRecordedPayload `json:"payload"`
}

type StockMovementRecordedPayload struct {
	MovementID    string    `json:"movement_id"`
	ProductID     string    `json:"product_id"`
	VariantID     string    `json:"variant_id,omitempty"`
	WarehouseID   string    `json:"warehouse_id"`
	SKU           string    `json:"sku"`
	MovementType  string    `json:"movement_type"`
	Quantity      int       `json:"quantity"`
	PreviousQty   int       `json:"previous_quantity"`
	NewQty        int       `json:"new_quantity"`
	Reference     string    `json:"reference,omitempty"`
	Reason        string    `json:"reason,omitempty"`
	PerformedBy   string    `json:"performed_by,omitempty"`
	RecordedAt    time.Time `json:"recorded_at"`
}

// InventorySnapshotEvent is published periodically for reporting/analytics
type InventorySnapshotEvent struct {
	BaseEvent
	Payload InventorySnapshotPayload `json:"payload"`
}

type InventorySnapshotPayload struct {
	SnapshotID    string              `json:"snapshot_id"`
	WarehouseID   string              `json:"warehouse_id,omitempty"`
	Items         []SnapshotItem      `json:"items"`
	TotalProducts int                 `json:"total_products"`
	TotalQuantity int                 `json:"total_quantity"`
	TotalReserved int                 `json:"total_reserved"`
	GeneratedAt   time.Time           `json:"generated_at"`
}

type SnapshotItem struct {
	ProductID       string `json:"product_id"`
	VariantID       string `json:"variant_id,omitempty"`
	SKU             string `json:"sku"`
	WarehouseID     string `json:"warehouse_id"`
	Quantity        int    `json:"quantity"`
	ReservedQty     int    `json:"reserved_quantity"`
	AvailableQty    int    `json:"available_quantity"`
}

// Serialize converts event to JSON bytes
func (e *StockReservedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockReservationFailedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockReleasedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockDecrementedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockReplenishedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *LowStockAlertEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *StockMovementRecordedEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *InventorySnapshotEvent) Serialize() ([]byte, error) {
	return json.Marshal(e)
}