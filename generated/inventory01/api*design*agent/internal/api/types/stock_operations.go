// file: internal/api/types/stock_operations.go
package types

import (
	"time"

	"github.com/google/uuid"
)

// ReservationItem represents a single item in a reservation request
type ReservationItem struct {
	// ProductID is the product to reserve
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	// Quantity to reserve
	Quantity int `json:"quantity" validate:"required,min=1"`
	// PreferredWarehouseID is optional; system will select if not provided
	PreferredWarehouseID *uuid.UUID `json:"preferred_warehouse_id,omitempty"`
}

// ReserveStockRequest is the request body for reserving stock
// @Description Request to reserve stock for an order
type ReserveStockRequest struct {
	// OrderID is the associated order identifier
	OrderID string `json:"order_id" validate:"required,min=1,max=100"`
	// Items contains the products and quantities to reserve
	Items []ReservationItem `json:"items" validate:"required,min=1,dive"`
	// ExpiresAt is when the reservation expires (optional, defaults to 30 min)
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// IdempotencyKey prevents duplicate reservations
	IdempotencyKey string `json:"idempotency_key" validate:"required,min=1,max=100"`
}

// ReservationResult represents the result of a single item reservation
type ReservationResult struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	Quantity      int       `json:"quantity"`
	ExpiresAt     Timestamp `json:"expires_at"`
}

// ReserveStockResponse is the response for stock reservation
type ReserveStockResponse struct {
	OrderID      string              `json:"order_id"`
	Reservations []ReservationResult `json:"reservations"`
	// PartialSuccess indicates some items couldn't be fully reserved
	PartialSuccess bool `json:"partial_success"`
	// FailedItems lists items that couldn't be reserved
	FailedItems []ReservationFailure `json:"failed_items,omitempty"`
}

// ReservationFailure describes why a reservation failed
type ReservationFailure struct {
	ProductID         uuid.UUID `json:"product_id"`
	RequestedQuantity int       `json:"requested_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	Reason            string    `json:"reason"`
}

// ReleaseStockRequest is the request body for releasing reserved stock
// @Description Request to release previously reserved stock
type ReleaseStockRequest struct {
	// OrderID releases all reservations for this order
	OrderID string `json:"order_id,omitempty" validate:"required_without=ReservationIDs"`
	// ReservationIDs releases specific reservations
	ReservationIDs []uuid.UUID `json:"reservation_ids,omitempty" validate:"required_without=OrderID"`
	// Reason for releasing the reservation
	Reason string `json:"reason" validate:"required,min=5,max=500"`
}

// ReleaseStockResponse is the response for stock release
type ReleaseStockResponse struct {
	ReleasedCount int       `json:"released_count"`
	ReleasedIDs   []uuid.UUID `json:"released_ids"`
}

// FulfillStockRequest is the request body for fulfillment
// @Description Request to decrement stock upon order fulfillment
type FulfillStockRequest struct {
	// OrderID identifies the order being fulfilled
	OrderID string `json:"order_id" validate:"required"`
	// ReservationIDs are the reservations to fulfill
	ReservationIDs []uuid.UUID `json:"reservation_ids" validate:"required,min=1"`
	// ShipmentID is the associated shipment identifier
	ShipmentID string `json:"shipment_id,omitempty"`
}

// FulfillStockResponse is the response for fulfillment
type FulfillStockResponse struct {
	FulfilledCount int                `json:"fulfilled_count"`
	Movements      []StockMovementRef `json:"movements"`
}

// StockMovementRef is a reference to a created stock movement
type StockMovementRef struct {
	MovementID  uuid.UUID `json:"movement_id"`
	ProductID   uuid.UUID `json:"product_id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
}

// ReplenishStockRequest is the request body for stock replenishment
// @Description Request to add stock from a supplier delivery
type ReplenishStockRequest struct {
	// WarehouseID is the receiving warehouse
	WarehouseID uuid.UUID `json:"warehouse_id" validate:"required"`
	// Items contains products and quantities to add
	Items []ReplenishmentItem `json:"items" validate:"required,min=1,dive"`
	// PurchaseOrderID is the associated PO number
	PurchaseOrderID string `json:"purchase_order_id,omitempty" validate:"max=100"`
	// SupplierID is the supplier identifier
	SupplierID string `json:"supplier_id,omitempty" validate:"max=100"`
	// Notes for the replenishment
	Notes string `json:"notes,omitempty" validate:"max=1000"`
}

// ReplenishmentItem represents a single item in a replenishment
type ReplenishmentItem struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
	// UnitCost is optional cost tracking
	UnitCost *float64 `json:"unit_cost,omitempty" validate:"omitempty,min=0"`
}

// ReplenishStockResponse is the response for replenishment
type ReplenishStockResponse struct {
	ReplenishedCount int                `json:"replenished_count"`
	Movements        []StockMovementRef `json:"movements"`
	AlertsCleared    []uuid.UUID        `json:"alerts_cleared,omitempty"`
}

// TransferStockRequest is the request body for warehouse transfers
// @Description Request to transfer stock between warehouses
type TransferStockRequest struct {
	// SourceWarehouseID is the origin warehouse
	SourceWarehouseID uuid.UUID `json:"source_warehouse_id" validate:"required"`
	// DestinationWarehouseID is the target warehouse
	DestinationWarehouseID uuid.UUID `json:"destination_warehouse_id" validate:"required,nefield=SourceWarehouseID"`
	// Items to transfer
	Items []TransferItem `json:"items" validate:"required,min=1,dive"`
	// TransferReference is an optional reference number
	TransferReference string `json:"transfer_reference,omitempty" validate:"max=100"`
	// Notes for the transfer
	Notes string `json:"notes,omitempty" validate:"max=1000"`
}

// TransferItem represents a single item in a transfer
type TransferItem struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// TransferStockResponse is the response for transfers
type TransferStockResponse struct {
	TransferID      uuid.UUID          `json:"transfer_id"`
	TransferredCount int               `json:"transferred_count"`
	Movements       []StockMovementRef `json:"movements"`
}