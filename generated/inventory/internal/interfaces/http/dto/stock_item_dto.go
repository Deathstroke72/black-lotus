// file: internal/interfaces/http/dto/stock_item_dto.go
package dto

import "time"

// CreateStockItemRequest represents the request body for creating a stock item.
// @Description Request payload for creating a stock item (product in warehouse)
type CreateStockItemRequest struct {
	// ProductID is the ID of the product
	ProductID string `json:"product_id" validate:"required,uuid"`
	// WarehouseID is the ID of the warehouse
	WarehouseID string `json:"warehouse_id" validate:"required,uuid"`
	// VariantSKU is the SKU of the specific variant (optional)
	VariantSKU string `json:"variant_sku,omitempty" validate:"max=100"`
	// Quantity is the initial stock quantity
	Quantity int `json:"quantity" validate:"min=0"`
	// ReorderPoint is the quantity at which to trigger reorder
	ReorderPoint int `json:"reorder_point" validate:"min=0"`
	// ReorderQuantity is the quantity to order when reordering
	ReorderQuantity int `json:"reorder_quantity" validate:"min=0"`
	// BinLocation is the physical location within the warehouse
	BinLocation string `json:"bin_location,omitempty" validate:"max=100"`
}

// StockItemResponse represents a stock item in API responses.
// @Description Stock item information returned by the API
type StockItemResponse struct {
	// ID is the unique stock item identifier
	ID string `json:"id"`
	// ProductID is the ID of the product
	ProductID string `json:"product_id"`
	// ProductName is the name of the product
	ProductName string `json:"product_name"`
	// WarehouseID is the ID of the warehouse
	WarehouseID string `json:"warehouse_id"`
	// WarehouseName is the name of the warehouse
	WarehouseName string `json:"warehouse_name"`
	// VariantSKU is the SKU of the specific variant
	VariantSKU string `json:"variant_sku,omitempty"`
	// Quantity is the current stock quantity
	Quantity int `json:"quantity"`
	// ReservedQuantity is the quantity currently reserved
	ReservedQuantity int `json:"reserved_quantity"`
	// AvailableQuantity is Quantity minus ReservedQuantity
	AvailableQuantity int `json:"available_quantity"`
	// ReorderPoint is the quantity at which to trigger reorder
	ReorderPoint int `json:"reorder_point"`
	// ReorderQuantity is the quantity to order when reordering
	ReorderQuantity int `json:"reorder_quantity"`
	// BinLocation is the physical location within the warehouse
	BinLocation string `json:"bin_location,omitempty"`
	// IsLowStock indicates if current quantity is below threshold
	IsLowStock bool `json:"is_low_stock"`
	// CreatedAt is when the stock item was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the stock item was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ListStockItemsRequest represents query parameters for listing stock items.
type ListStockItemsRequest struct {
	PaginationRequest
	// ProductID filters by product
	ProductID string `json:"product_id,omitempty" validate:"omitempty,uuid"`
	// WarehouseID filters by warehouse
	WarehouseID string `json:"warehouse_id,omitempty" validate:"omitempty,uuid"`
	// VariantSKU filters by variant SKU
	VariantSKU string `json:"variant_sku,omitempty"`
	// LowStockOnly returns only items with low stock
	LowStockOnly bool `json:"low_stock_only,omitempty"`
}

// ListStockItemsResponse represents the response for listing stock items.
// @Description Paginated list of stock items
type ListStockItemsResponse struct {
	// StockItems is the list of stock items
	StockItems []StockItemResponse `json:"stock_items"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}

// AggregatedStockResponse represents aggregated stock across warehouses.
// @Description Aggregated stock information for a product
type AggregatedStockResponse struct {
	// ProductID is the product identifier
	ProductID string `json:"product_id"`
	// ProductName is the product name
	ProductName string `json:"product_name"`
	// TotalQuantity is the total stock across all warehouses
	TotalQuantity int `json:"total_quantity"`
	// TotalReserved is the total reserved quantity
	TotalReserved int `json:"total_reserved"`
	// TotalAvailable is TotalQuantity minus TotalReserved
	TotalAvailable int `json:"total_available"`
	// IsLowStock indicates if total stock is below threshold
	IsLowStock bool `json:"is_low_stock"`
	// WarehouseBreakdown shows stock per warehouse
	WarehouseBreakdown []WarehouseStockBreakdown `json:"warehouse_breakdown"`
	// VariantBreakdown shows stock per variant
	VariantBreakdown []VariantStockBreakdown `json:"variant_breakdown,omitempty"`
}

// WarehouseStockBreakdown shows stock for a specific warehouse.
type WarehouseStockBreakdown struct {
	// WarehouseID is the warehouse identifier
	WarehouseID string `json:"warehouse_id"`
	// WarehouseName is the warehouse name
	WarehouseName string `json:"warehouse_name"`
	// Quantity is the stock in this warehouse
	Quantity int `json:"quantity"`
	// Reserved is the reserved quantity in this warehouse
	Reserved int `json:"reserved"`
	// Available is the available quantity
	Available int `json:"available"`
}

// VariantStockBreakdown shows stock for a specific variant.
type VariantStockBreakdown struct {
	// VariantSKU is the variant SKU
	VariantSKU string `json:"variant_sku"`
	// Size is the variant size
	Size string `json:"size,omitempty"`
	// Color is the variant color
	Color string `json:"color,omitempty"`
	// Quantity is the total stock for this variant
	Quantity int `json:"quantity"`
	// Reserved is the reserved quantity for this variant
	Reserved int `json:"reserved"`
	// Available is the available quantity
	Available int `json:"available"`
}