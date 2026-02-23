// file: internal/api/types/movement.go
package types

import (
	"time"

	"github.com/google/uuid"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeReservation  MovementType = "RESERVATION"
	MovementTypeRelease      MovementType = "RELEASE"
	MovementTypeFulfillment  MovementType = "FULFILLMENT"
	MovementTypeReplenishment MovementType = "REPLENISHMENT"
	MovementTypeTransferOut  MovementType = "TRANSFER_OUT"
	MovementTypeTransferIn   MovementType = "TRANSFER_IN"
	MovementTypeAdjustment   MovementType = "ADJUSTMENT"
)

// StockMovement represents an audit record of stock changes
type StockMovement struct {
	ID               uuid.UUID    `json:"id"`
	ProductID        uuid.UUID    `json:"product_id"`
	ProductSKU       string       `json:"product_sku"`
	WarehouseID      uuid.UUID    `json:"warehouse_id"`
	WarehouseCode    string       `json:"warehouse_code"`
	MovementType     MovementType `json:"movement_type"`
	Quantity         int          `json:"quantity"` // Positive for additions, negative for deductions
	QuantityBefore   int          `json:"quantity_before"`
	QuantityAfter    int          `json:"quantity_after"`
	ReferenceType    string       `json:"reference_type,omitempty"` // ORDER, TRANSFER, PO, etc.
	ReferenceID      string       `json:"reference_id,omitempty"`
	ReservationID    *uuid.UUID   `json:"reservation_id,omitempty"`
	Notes            string       `json:"notes,omitempty"`
	CreatedAt        Timestamp    `json:"created_at"`
	CreatedBy        string       `json:"created_by"`
}

// MovementQueryParams contains query parameters for movement queries
type MovementQueryParams struct {
	Pagination
	ProductID     *uuid.UUID    `json:"product_id,omitempty"`
	WarehouseID   *uuid.UUID    `json:"warehouse_id,omitempty"`
	MovementType  *MovementType `json:"movement_type,omitempty"`
	ReferenceType *string       `json:"reference_type,omitempty"`
	ReferenceID   *string       `json:"reference_id,omitempty"`
	StartDate     *time.Time    `json:"start_date,omitempty"`
	EndDate       *time.Time    `json:"end_date,omitempty"`
}

// MovementResponse wraps single movement responses
type MovementResponse struct {
	Movement StockMovement `json:"movement"`
}

// MovementListResponse is the response for listing movements
type MovementListResponse = PaginatedResponse[StockMovement]