# API Design Agent Output

# Inventory Service REST API Design

## API Endpoints Overview

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/products` | Create a new product |
| GET | `/api/v1/products` | List products with pagination |
| GET | `/api/v1/products/{productId}` | Get product by ID |
| PUT | `/api/v1/products/{productId}` | Update product details |
| DELETE | `/api/v1/products/{productId}` | Soft delete a product |
| POST | `/api/v1/warehouses` | Create a new warehouse |
| GET | `/api/v1/warehouses` | List all warehouses |
| GET | `/api/v1/warehouses/{warehouseId}` | Get warehouse by ID |
| PUT | `/api/v1/warehouses/{warehouseId}` | Update warehouse details |
| DELETE | `/api/v1/warehouses/{warehouseId}` | Soft delete a warehouse |
| POST | `/api/v1/stock-items` | Create stock item (product in warehouse) |
| GET | `/api/v1/stock-items` | List stock items with filters |
| GET | `/api/v1/stock-items/{stockItemId}` | Get stock item by ID |
| GET | `/api/v1/products/{productId}/stock` | Get aggregated stock across warehouses |
| POST | `/api/v1/reservations` | Reserve stock for an order |
| GET | `/api/v1/reservations/{reservationId}` | Get reservation details |
| POST | `/api/v1/reservations/{reservationId}/release` | Release reserved stock |
| POST | `/api/v1/reservations/{reservationId}/fulfill` | Fulfill reservation (decrement stock) |
| GET | `/api/v1/orders/{orderId}/reservations` | Get reservations by order ID |
| POST | `/api/v1/stock-movements/replenish` | Replenish stock |
| GET | `/api/v1/stock-movements` | List stock movements (audit trail) |
| GET | `/api/v1/stock-items/{stockItemId}/movements` | Get movements for specific stock item |
| GET | `/api/v1/alerts/low-stock` | Get current low-stock alerts |

---

## DTOs

```go
// file: internal/interfaces/http/dto/error_dto.go
package dto

import "time"

// ErrorResponse represents the standardized error response envelope.
// @Description Standard error response format for all API errors
type ErrorResponse struct {
	// Error contains the error details
	Error ErrorDetail `json:"error"`
	// RequestID is the unique identifier for request tracing
	RequestID string `json:"request_id,omitempty"`
}

// ErrorDetail contains specific error information.
type ErrorDetail struct {
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Message is a human-readable error description
	Message string `json:"message"`
	// Details contains additional error context
	Details []FieldError `json:"details,omitempty"`
	// Timestamp is when the error occurred
	Timestamp time.Time `json:"timestamp"`
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	// Field is the name of the field that failed validation
	Field string `json:"field"`
	// Message describes why validation failed
	Message string `json:"message"`
}

// Common error codes
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInsufficientStock = "INSUFFICIENT_STOCK"
	ErrCodeInvalidState     = "INVALID_STATE"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
)
```

```go
// file: internal/interfaces/http/dto/pagination_dto.go
package dto

// PaginationRequest represents pagination parameters for list endpoints.
type PaginationRequest struct {
	// Page number (1-indexed)
	Page int `json:"page" validate:"min=1"`
	// PageSize is the number of items per page
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse contains pagination metadata in list responses.
type PaginationResponse struct {
	// Page is the current page number
	Page int `json:"page"`
	// PageSize is the number of items per page
	PageSize int `json:"page_size"`
	// TotalItems is the total number of items across all pages
	TotalItems int64 `json:"total_items"`
	// TotalPages is the total number of pages
	TotalPages int `json:"total_pages"`
	// HasNext indicates if there are more pages
	HasNext bool `json:"has_next"`
	// HasPrev indicates if there are previous pages
	HasPrev bool `json:"has_prev"`
}

// DefaultPage is the default page number
const DefaultPage = 1

// DefaultPageSize is the default number of items per page
const DefaultPageSize = 20

// MaxPageSize is the maximum allowed page size
const MaxPageSize = 100
```

```go
// file: internal/interfaces/http/dto/product_dto.go
package dto

import "time"

// ProductVariant represents a product variant (size, color combination).
type ProductVariant struct {
	// Size of the product variant (e.g., "S", "M", "L", "XL")
	Size string `json:"size,omitempty"`
	// Color of the product variant (e.g., "Red", "Blue")
	Color string `json:"color,omitempty"`
	// SKU is the unique stock keeping unit for this variant
	SKU string `json:"sku" validate:"required,min=1,max=100"`
}

// CreateProductRequest represents the request body for creating a product.
// @Description Request payload for creating a new product
type CreateProductRequest struct {
	// Name is the product display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Description is the product description
	Description string `json:"description,omitempty" validate:"max=2000"`
	// BaseSKU is the base SKU for the product (variants will extend this)
	BaseSKU string `json:"base_sku" validate:"required,min=1,max=100"`
	// Category is the product category
	Category string `json:"category,omitempty" validate:"max=100"`
	// Variants are the product variants (size, color combinations)
	Variants []ProductVariant `json:"variants,omitempty" validate:"dive"`
	// LowStockThreshold is the quantity below which low-stock alerts trigger
	LowStockThreshold int `json:"low_stock_threshold" validate:"min=0"`
	// Metadata contains additional product attributes
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateProductRequest represents the request body for updating a product.
// @Description Request payload for updating an existing product
type UpdateProductRequest struct {
	// Name is the product display name
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	// Description is the product description
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	// Category is the product category
	Category *string `json:"category,omitempty" validate:"omitempty,max=100"`
	// LowStockThreshold is the quantity below which low-stock alerts trigger
	LowStockThreshold *int `json:"low_stock_threshold,omitempty" validate:"omitempty,min=0"`
	// Metadata contains additional product attributes
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ProductResponse represents a product in API responses.
// @Description Product information returned by the API
type ProductResponse struct {
	// ID is the unique product identifier
	ID string `json:"id"`
	// Name is the product display name
	Name string `json:"name"`
	// Description is the product description
	Description string `json:"description,omitempty"`
	// BaseSKU is the base SKU for the product
	BaseSKU string `json:"base_sku"`
	// Category is the product category
	Category string `json:"category,omitempty"`
	// Variants are the product variants
	Variants []ProductVariant `json:"variants,omitempty"`
	// LowStockThreshold is the quantity below which low-stock alerts trigger
	LowStockThreshold int `json:"low_stock_threshold"`
	// TotalStock is the aggregated stock across all warehouses
	TotalStock int `json:"total_stock"`
	// TotalReserved is the aggregated reserved quantity
	TotalReserved int `json:"total_reserved"`
	// AvailableStock is TotalStock minus TotalReserved
	AvailableStock int `json:"available_stock"`
	// Metadata contains additional product attributes
	Metadata map[string]string `json:"metadata,omitempty"`
	// CreatedAt is when the product was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the product was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ListProductsRequest represents query parameters for listing products.
type ListProductsRequest struct {
	PaginationRequest
	// Category filters by product category
	Category string `json:"category,omitempty"`
	// Search performs a text search on name and description
	Search string `json:"search,omitempty"`
	// LowStockOnly returns only products with low stock
	LowStockOnly bool `json:"low_stock_only,omitempty"`
}

// ListProductsResponse represents the response for listing products.
// @Description Paginated list of products
type ListProductsResponse struct {
	// Products is the list of products
	Products []ProductResponse `json:"products"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}
```

```go
// file: internal/interfaces/http/dto/warehouse_dto.go
package dto

import "time"

// WarehouseAddress represents a warehouse physical address.
type WarehouseAddress struct {
	// Street is the street address
	Street string `json:"street" validate:"required,max=255"`
	// City is the city name
	City string `json:"city" validate:"required,max=100"`
	// State is the state or province
	State string `json:"state,omitempty" validate:"max=100"`
	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code" validate:"required,max=20"`
	// Country is the ISO 3166-1 alpha-2 country code
	Country string `json:"country" validate:"required,len=2"`
}

// CreateWarehouseRequest represents the request body for creating a warehouse.
// @Description Request payload for creating a new warehouse
type CreateWarehouseRequest struct {
	// Name is the warehouse display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Code is a unique warehouse code (e.g., "WH-NYC-01")
	Code string `json:"code" validate:"required,min=1,max=50"`
	// Address is the warehouse physical address
	Address WarehouseAddress `json:"address" validate:"required"`
	// IsActive indicates if the warehouse is operational
	IsActive bool `json:"is_active"`
	// Priority is used for stock allocation (lower = higher priority)
	Priority int `json:"priority" validate:"min=0"`
}

// UpdateWarehouseRequest represents the request body for updating a warehouse.
// @Description Request payload for updating an existing warehouse
type UpdateWarehouseRequest struct {
	// Name is the warehouse display name
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	// Address is the warehouse physical address
	Address *WarehouseAddress `json:"address,omitempty"`
	// IsActive indicates if the warehouse is operational
	IsActive *bool `json:"is_active,omitempty"`
	// Priority is used for stock allocation
	Priority *int `json:"priority,omitempty" validate:"omitempty,min=0"`
}

// WarehouseResponse represents a warehouse in API responses.
// @Description Warehouse information returned by the API
type WarehouseResponse struct {
	// ID is the unique warehouse identifier
	ID string `json:"id"`
	// Name is the warehouse display name
	Name string `json:"name"`
	// Code is the unique warehouse code
	Code string `json:"code"`
	// Address is the warehouse physical address
	Address WarehouseAddress `json:"address"`
	// IsActive indicates if the warehouse is operational
	IsActive bool `json:"is_active"`
	// Priority is used for stock allocation
	Priority int `json:"priority"`
	// TotalProducts is the count of distinct products in this warehouse
	TotalProducts int `json:"total_products"`
	// CreatedAt is when the warehouse was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the warehouse was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ListWarehousesRequest represents query parameters for listing warehouses.
type ListWarehousesRequest struct {
	PaginationRequest
	// ActiveOnly returns only active warehouses
	ActiveOnly bool `json:"active_only,omitempty"`
}

// ListWarehousesResponse represents the response for listing warehouses.
// @Description List of warehouses
type ListWarehousesResponse struct {
	// Warehouses is the list of warehouses
	Warehouses []WarehouseResponse `json:"warehouses"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}
```

```go
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
```

```go
// file: internal/interfaces/http/dto/reservation_dto.go
package dto

import "time"

// ReservationItem represents a single item in a reservation request.
type ReservationItem struct {
	// ProductID is the product to reserve
	ProductID string `json:"product_id" validate:"required,uuid"`
	// VariantSKU is the specific variant SKU (optional)
	VariantSKU string `json:"variant_sku,omitempty" validate:"max=100"`
	// Quantity is the amount to reserve
	Quantity int `json:"quantity" validate:"required,min=1"`
	// PreferredWarehouseID is the preferred warehouse (optional)
	PreferredWarehouseID string `json:"preferred_warehouse_id,omitempty" validate:"omitempty,uuid"`
}

// CreateReservationRequest represents the request body for creating a reservation.
// @Description Request payload for reserving stock for an order
type CreateReservationRequest struct {
	// OrderID is the external order identifier
	OrderID string `json:"order_id" validate:"required,min=1,max=100"`
	// Items are the products and quantities to reserve
	Items []ReservationItem `json:"items" validate:"required,min=1,dive"`
	// ExpiresAt is when the reservation should expire (optional)
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Metadata contains additional reservation context
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ReservationItemResponse represents a reserved item in the response.
type ReservationItemResponse struct {
	// ProductID is the reserved product
	ProductID string `json:"product_id"`
	// ProductName is the product name
	ProductName string `json:"product_name"`
	// VariantSKU is the variant SKU
	VariantSKU string `json:"variant_sku,omitempty"`
	// Quantity is the reserved amount
	Quantity int `json:"quantity"`
	// WarehouseID is where the stock is reserved
	WarehouseID string `json:"warehouse_id"`
	// WarehouseName is the warehouse name
	WarehouseName string `json:"warehouse_name"`
	// StockItemID is the specific stock item
	StockItemID string `json:"stock_item_id"`
}

// ReservationResponse represents a reservation in API responses.
// @Description Reservation information returned by the API
type ReservationResponse struct {
	// ID is the unique reservation identifier
	ID string `json:"id"`
	// OrderID is the external order identifier
	OrderID string `json:"order_id"`
	// Status is the reservation status (pending, confirmed, released, fulfilled, expired)
	Status string `json:"status"`
	// Items are the reserved items
	Items []ReservationItemResponse `json:"items"`
	// ExpiresAt is when the reservation expires
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Metadata contains additional reservation context
	Metadata map[string]string `json:"metadata,omitempty"`
	// CreatedAt is when the reservation was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the reservation was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ReleaseReservationRequest represents the request body for releasing a reservation.
// @Description Request payload for releasing reserved stock
type ReleaseReservationRequest struct {
	// Reason is the reason for releasing (e.g., "order_cancelled", "timeout")
	Reason string `json:"reason" validate:"required,min=1,max=255"`
	// PartialItems allows releasing only specific items (optional)
	PartialItems []PartialReleaseItem `json:"partial_items,omitempty" validate:"omitempty,dive"`
}

// PartialReleaseItem specifies partial release of a reservation item.
type PartialReleaseItem struct {
	// ProductID is the product to partially release
	ProductID string `json:"product_id" validate:"required,uuid"`
	// VariantSKU is the variant SKU
	VariantSKU string `json:"variant_sku,omitempty"`
	// Quantity is the amount to release
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// FulfillReservationRequest represents the request body for fulfilling a reservation.
// @Description Request payload for fulfilling a reservation (decrementing stock)
type FulfillReservationRequest struct {
	// ShipmentID is the external shipment identifier (optional)
	ShipmentID string `json:"shipment_id,omitempty" validate:"max=100"`
	// FulfilledBy is the user or system that fulfilled the reservation
	FulfilledBy string `json:"fulfilled_by" validate:"required,max=255"`
	// Notes contains any fulfillment notes
	Notes string `json:"notes,omitempty" validate:"max=1000"`
}

// ListReservationsRequest represents query parameters for listing reservations.
type ListReservationsRequest struct {
	PaginationRequest
	// OrderID filters by order
	OrderID string `json:"order_id,omitempty"`
	// Status filters by reservation status
	Status string `json:"status,omitempty"`
	// ProductID filters by product
	ProductID string `json:"product_id,omitempty" validate:"omitempty,uuid"`
}

// ListReservationsResponse represents the response for listing reservations.
// @Description Paginated list of reservations
type ListReservationsResponse struct {
	// Reservations is the list of reservations
	Reservations []ReservationResponse `json:"reservations"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}

// Reservation status constants
const (
	ReservationStatusPending   = "pending"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusReleased  = "released"
	ReservationStatusFulfilled = "fulfilled"
	ReservationStatusExpired   = "expired"
)
```

```go
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
```

```go
// file: internal/interfaces/http/dto/alert_dto.go
package dto

import "time"

// LowStockAlertResponse represents a low stock alert.
// @Description Low stock alert information
type LowStockAlertResponse struct {
	// ID is the unique alert identifier
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
	// CurrentQuantity is the current stock level
	CurrentQuantity int `json:"current_quantity"`
	// Threshold is the low stock threshold
	Threshold int `json:"threshold"`
	// ReorderPoint