// file: internal/interfaces/http/dto/stock_movement_dto.go
package dto

import "time"

// ReplenishStockRequest represents the request body for replenishing stock.
// @Description Request payload for replenishing stock in a warehouse
type ReplenishStockRequest struct {
	// StockItemID is the stock item to replenish
	StockItemID string `json:"stock_item_id" validate:"required,uuid"`
	// Quantity is the amount to add
	Quantity int `json:"quantity" validate:"required,min=1"`
	// ReferenceType is the type of reference (e.g., "purchase_order", "transfer", "adjustment")
	ReferenceType string `json:"reference_type" validate:"required,oneof=purchase_order transfer adjustment return"`
	// ReferenceID is the external reference identifier
	ReferenceID string `json:"reference_id" validate:"required,min=1,max=100"`
	// UnitCost is the cost per unit (optional)
	UnitCost *float64 `json:"unit_cost,omitempty" validate:"omitempty,min=0"`
	// Notes contains any additional notes
	Notes string `json:"notes,omitempty" validate:"max=1000"`
	// PerformedBy is the user who performed the replenishment
	PerformedBy string `json:"performed_by" validate:"required,max=255"`
}

// StockMovementResponse represents a stock movement in API responses.
// @Description Stock movement information for audit trail
type StockMovementResponse struct {
	// ID is the unique movement identifier
	ID string `json:"id"`
	// StockItemID is the affected stock item
	StockItemID string `json:"stock_item_id"`
	// ProductID is the product identifier
	ProductID string `json:"product_id"`
	// ProductName is the product name
	ProductName string `json:"product_name"`
	// WarehouseID is the warehouse identifier
	WarehouseID string `json:"warehouse_id"`
	// WarehouseName is the warehouse name
	WarehouseName string `json:"warehouse_name"`
	// VariantSKU is the variant SKU
	VariantSKU string `json:"variant_sku,omitempty"`
	// MovementType is the type of movement
	MovementType string `json:"movement_type"`
	// Quantity is the quantity changed (positive for in, negative for out)
	Quantity int `json:"quantity"`
	// QuantityBefore is the quantity before the movement
	QuantityBefore int `json:"quantity_before"`
	// QuantityAfter is the quantity after the movement
	QuantityAfter int `json:"quantity_after"`
	// ReferenceType is the type of reference
	ReferenceType string `json:"reference_type"`
	// ReferenceID is the external reference identifier
	ReferenceID string `json:"reference_id"`
	// UnitCost is the cost per unit
	UnitCost *float64 `json:"unit_cost,omitempty"`
	// Notes contains additional notes
	Notes string `json:"notes,omitempty"`
	// PerformedBy is who performed the movement
	PerformedBy string `json:"performed_by"`
	// CreatedAt is when the movement occurred
	CreatedAt time.Time `json:"created_at"`
}

// ListStockMovementsRequest represents query parameters for listing movements.
type ListStockMovementsRequest struct {
	PaginationRequest
	// StockItemID filters by stock item
	StockItemID string `json:"stock_item_id,omitempty" validate:"omitempty,uuid"`
	// ProductID filters by product
	ProductID string `json:"product_id,omitempty" validate:"omitempty,uuid"`
	// WarehouseID filters by warehouse
	WarehouseID string `json:"warehouse_id,omitempty" validate:"omitempty,uuid"`
	// MovementType filters by movement type
	MovementType string `json:"movement_type,omitempty"`
	// ReferenceType filters by reference type
	ReferenceType string `json:"reference_type,omitempty"`
	// ReferenceID filters by reference ID
	ReferenceID string `json:"reference_id,omitempty"`
	// StartDate filters movements from this date
	StartDate *time.Time `json:"start_date,omitempty"`
	// EndDate filters movements until this date
	EndDate *time.Time `json:"end_date,omitempty"`
}

// ListStockMovementsResponse represents the response for listing movements.
// @Description Paginated list of stock movements (audit trail)
type ListStockMovementsResponse struct {
	// Movements is the list of stock movements
	Movements []StockMovementResponse `json:"movements"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}

// Movement type constants
const (
	MovementTypeReplenish   = "replenish"
	MovementTypeReserve     = "reserve"
	MovementTypeRelease     = "release"
	MovementTypeFulfill     = "fulfill"
	MovementTypeAdjustment  = "adjustment"
	MovementTypeTransferIn  = "transfer_in"
	MovementTypeTransferOut = "transfer_out"
)