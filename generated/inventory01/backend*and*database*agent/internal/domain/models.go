// file: internal/domain/models.go
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Product represents a product in the inventory system
type Product struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	SKU               string     `json:"sku" db:"sku"`
	Name              string     `json:"name" db:"name"`
	Description       *string    `json:"description,omitempty" db:"description"`
	Category          *string    `json:"category,omitempty" db:"category"`
	BasePrice         *float64   `json:"base_price,omitempty" db:"base_price"`
	LowStockThreshold int        `json:"low_stock_threshold" db:"low_stock_threshold"`
	IsActive          bool       `json:"is_active" db:"is_active"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ProductVariant represents a variant of a product (size, color, etc.)
type ProductVariant struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	ProductID         uuid.UUID       `json:"product_id" db:"product_id"`
	SKU               string          `json:"sku" db:"sku"`
	Name              string          `json:"name" db:"name"`
	Attributes        json.RawMessage `json:"attributes" db:"attributes"`
	PriceModifier     *float64        `json:"price_modifier,omitempty" db:"price_modifier"`
	LowStockThreshold *int            `json:"low_stock_threshold,omitempty" db:"low_stock_threshold"`
	IsActive          bool            `json:"is_active" db:"is_active"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
}

// VariantAttributes represents decoded variant attributes
type VariantAttributes struct {
	Size   *string `json:"size,omitempty"`
	Color  *string `json:"color,omitempty"`
	Weight *string `json:"weight,omitempty"`
	Custom map[string]string `json:"custom,omitempty"`
}

// Warehouse represents a storage location
type Warehouse struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Code      string          `json:"code" db:"code"`
	Name      string          `json:"name" db:"name"`
	Address   json.RawMessage `json:"address" db:"address"`
	IsActive  bool            `json:"is_active" db:"is_active"`
	Priority  int             `json:"priority" db:"priority"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
}

// WarehouseAddress represents the address structure for a warehouse
type WarehouseAddress struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// StockItem represents inventory at a specific warehouse
type StockItem struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	ProductID        uuid.UUID  `json:"product_id" db:"product_id"`
	VariantID        *uuid.UUID `json:"variant_id,omitempty" db:"variant_id"`
	WarehouseID      uuid.UUID  `json:"warehouse_id" db:"warehouse_id"`
	Quantity         int        `json:"quantity" db:"quantity"`
	ReservedQuantity int        `json:"reserved_quantity" db:"reserved_quantity"`
	ReorderPoint     int        `json:"reorder_point" db:"reorder_point"`
	ReorderQuantity  int        `json:"reorder_quantity" db:"reorder_quantity"`
	LastCountedAt    *time.Time `json:"last_counted_at,omitempty" db:"last_counted_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	Version          int        `json:"version" db:"version"`
}

// AvailableQuantity returns the quantity available for reservation
func (s *StockItem) AvailableQuantity() int {
	return s.Quantity - s.ReservedQuantity
}

// StockMovementType represents the type of stock movement
type StockMovementType string

const (
	MovementTypeReservation  StockMovementType = "RESERVATION"
	MovementTypeRelease      StockMovementType = "RELEASE"
	MovementTypeFulfillment  StockMovementType = "FULFILLMENT"
	MovementTypeReplenishment StockMovementType = "REPLENISHMENT"
	MovementTypeAdjustment   StockMovementType = "ADJUSTMENT"
	MovementTypeTransferOut  StockMovementType = "TRANSFER_OUT"
	MovementTypeTransferIn   StockMovementType = "TRANSFER_IN"
	MovementTypeReturn       StockMovementType = "RETURN"
	MovementTypeDamage       StockMovementType = "DAMAGE"
	MovementTypeExpired      StockMovementType = "EXPIRED"
)

// StockMovement represents a change in stock level (audit trail)
type StockMovement struct {
	ID             uuid.UUID         `json:"id" db:"id"`
	StockItemID    uuid.UUID         `json:"stock_item_id" db:"stock_item_id"`
	ProductID      uuid.UUID         `json:"product_id" db:"product_id"`
	VariantID      *uuid.UUID        `json:"variant_id,omitempty" db:"variant_id"`
	WarehouseID    uuid.UUID         `json:"warehouse_id" db:"warehouse_id"`
	MovementType   StockMovementType `json:"movement_type" db:"movement_type"`
	Quantity       int               `json:"quantity" db:"quantity"`
	QuantityBefore int               `json:"quantity_before" db:"quantity_before"`
	QuantityAfter  int               `json:"quantity_after" db:"quantity_after"`
	ReservedBefore int               `json:"reserved_before" db:"reserved_before"`
	ReservedAfter  int               `json:"reserved_after" db:"reserved_after"`
	ReferenceType  *string           `json:"reference_type,omitempty" db:"reference_type"`
	ReferenceID    *uuid.UUID        `json:"reference_id,omitempty" db:"reference_id"`
	Reason         *string           `json:"reason,omitempty" db:"reason"`
	Metadata       json.RawMessage   `json:"metadata,omitempty" db:"metadata"`
	PerformedBy    *string           `json:"performed_by,omitempty" db:"performed_by"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
}

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusPending           ReservationStatus = "PENDING"
	ReservationStatusConfirmed         ReservationStatus = "CONFIRMED"
	ReservationStatusPartiallyFulfilled ReservationStatus = "PARTIALLY_FULFILLED"
	ReservationStatusFulfilled         ReservationStatus = "FULFILLED"
	ReservationStatusCancelled         ReservationStatus = "CANCELLED"
	ReservationStatusExpired           ReservationStatus = "EXPIRED"
)

// Reservation represents a stock reservation for an order
type Reservation struct {
	ID                uuid.UUID         `json:"id" db:"id"`
	OrderID           uuid.UUID         `json:"order_id" db:"order_id"`
	ProductID         uuid.UUID         `json:"product_id" db:"product_id"`
	VariantID         *uuid.UUID        `json:"variant_id,omitempty" db:"variant_id"`
	WarehouseID       uuid.UUID         `json:"warehouse_id" db:"warehouse_id"`
	StockItemID       uuid.UUID         `json:"stock_item_id" db:"stock_item_id"`
	Quantity          int               `json:"quantity" db:"quantity"`
	FulfilledQuantity int               `json:"fulfilled_quantity" db:"fulfilled_quantity"`
	Status            ReservationStatus `json:"status" db:"status"`
	ExpiresAt         time.Time         `json:"expires_at" db:"expires_at"`
	ConfirmedAt       *time.Time        `json:"confirmed_at,omitempty" db:"confirmed_at"`
	FulfilledAt       *time.Time        `json:"fulfilled_at,omitempty" db:"fulfilled_at"`
	CancelledAt       *time.Time        `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at" db:"updated_at"`
}

// RemainingQuantity returns the quantity not yet fulfilled
func (r *Reservation) RemainingQuantity() int {
	return r.Quantity - r.FulfilledQuantity
}

// LowStockAlert represents a low stock alert
type LowStockAlert struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ProductID       uuid.UUID  `json:"product_id" db:"product_id"`
	VariantID       *uuid.UUID `json:"variant_id,omitempty" db:"variant_id"`
	WarehouseID     *uuid.UUID `json:"warehouse_id,omitempty" db:"warehouse_id"`
	CurrentQuantity int        `json:"current_quantity" db:"current_quantity"`
	Threshold       int        `json:"threshold" db:"threshold"`
	IsAcknowledged  bool       `json:"is_acknowledged" db:"is_acknowledged"`
	AcknowledgedBy  *string    `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	AcknowledgedAt  *time.Time `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	ResolvedAt      *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
}

// OutboxEvent represents an event to be published
type OutboxEvent struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	AggregateType string          `json:"aggregate_type" db:"aggregate_type"`
	AggregateID   uuid.UUID       `json:"aggregate_id" db:"aggregate_id"`
	EventType     string          `json:"event_type" db:"event_type"`
	Payload       json.RawMessage `json:"payload" db:"payload"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	PublishedAt   *time.Time      `json:"published_at,omitempty" db:"published_at"`
	RetryCount    int             `json:"retry_count" db:"retry_count"`
	LastError     *string         `json:"last_error,omitempty" db:"last_error"`
}