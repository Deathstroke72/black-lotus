// file: internal/api/types/stock.go
package types

import (
	"github.com/google/uuid"
)

// StockItem represents the stock level of a product at a specific warehouse
type StockItem struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	WarehouseID     uuid.UUID `json:"warehouse_id"`
	QuantityOnHand  int       `json:"quantity_on_hand"`
	QuantityReserved int      `json:"quantity_reserved"`
	QuantityAvailable int     `json:"quantity_available"`
	Version         int64     `json:"version"` // For optimistic locking
	AuditInfo
}

// AggregatedStock represents stock aggregated across warehouses
type AggregatedStock struct {
	ProductID         uuid.UUID          `json:"product_id"`
	ProductSKU        string             `json:"product_sku"`
	ProductName       string             `json:"product_name"`
	TotalOnHand       int                `json:"total_on_hand"`
	TotalReserved     int                `json:"total_reserved"`
	TotalAvailable    int                `json:"total_available"`
	WarehouseBreakdown []WarehouseStock  `json:"warehouse_breakdown,omitempty"`
	LowStockAlert     bool               `json:"low_stock_alert"`
}

// WarehouseStock represents stock at a single warehouse
type WarehouseStock struct {
	WarehouseID       uuid.UUID `json:"warehouse_id"`
	WarehouseCode     string    `json:"warehouse_code"`
	WarehouseName     string    `json:"warehouse_name"`
	QuantityOnHand    int       `json:"quantity_on_hand"`
	QuantityReserved  int       `json:"quantity_reserved"`
	QuantityAvailable int       `json:"quantity_available"`
}

// StockQueryParams contains query parameters for stock queries
type StockQueryParams struct {
	Pagination
	ProductIDs    []uuid.UUID `json:"product_ids,omitempty"`
	WarehouseIDs  []uuid.UUID `json:"warehouse_ids,omitempty"`
	LowStockOnly  bool        `json:"low_stock_only,omitempty"`
	Category      string      `json:"category,omitempty"`
}

// SetStockRequest is used for admin stock level overrides
type SetStockRequest struct {
	// QuantityOnHand is the new stock quantity
	QuantityOnHand int `json:"quantity_on_hand" validate:"min=0"`
	// Reason for the adjustment (required for audit)
	Reason string `json:"reason" validate:"required,min=10,max=500"`
	// ExpectedVersion for optimistic locking (optional)
	ExpectedVersion *int64 `json:"expected_version,omitempty"`
}

// StockResponse wraps single stock item responses
type StockResponse struct {
	Stock StockItem `json:"stock"`
}

// AggregatedStockResponse wraps aggregated stock responses
type AggregatedStockResponse struct {
	Stock AggregatedStock `json:"stock"`
}

// StockListResponse is the response for listing stock items
type StockListResponse = PaginatedResponse[AggregatedStock]