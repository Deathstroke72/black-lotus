# API Design Agent Output

# Inventory Microservice API Design

## 1. Complete List of API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| **Products** |||
| GET | `/api/v1/products` | List all products with pagination |
| POST | `/api/v1/products` | Create a new product |
| GET | `/api/v1/products/{productId}` | Get product details |
| PUT | `/api/v1/products/{productId}` | Update product information |
| DELETE | `/api/v1/products/{productId}` | Soft-delete a product |
| GET | `/api/v1/products/{productId}/variants` | List all variants for a product |
| POST | `/api/v1/products/{productId}/variants` | Create a product variant |
| **Warehouses** |||
| GET | `/api/v1/warehouses` | List all warehouses |
| POST | `/api/v1/warehouses` | Create a new warehouse |
| GET | `/api/v1/warehouses/{warehouseId}` | Get warehouse details |
| PUT | `/api/v1/warehouses/{warehouseId}` | Update warehouse information |
| DELETE | `/api/v1/warehouses/{warehouseId}` | Deactivate a warehouse |
| **Stock Items** |||
| GET | `/api/v1/stock` | List stock items with filters |
| GET | `/api/v1/stock/{productId}` | Get aggregated stock across all warehouses |
| GET | `/api/v1/stock/{productId}/warehouses/{warehouseId}` | Get stock for product in specific warehouse |
| PUT | `/api/v1/stock/{productId}/warehouses/{warehouseId}` | Set stock level (admin override) |
| **Stock Operations** |||
| POST | `/api/v1/stock/reserve` | Reserve stock for an order |
| POST | `/api/v1/stock/release` | Release reserved stock |
| POST | `/api/v1/stock/fulfill` | Decrement stock on fulfillment |
| POST | `/api/v1/stock/replenish` | Replenish stock |
| POST | `/api/v1/stock/transfer` | Transfer stock between warehouses |
| **Reservations** |||
| GET | `/api/v1/reservations` | List reservations with filters |
| GET | `/api/v1/reservations/{reservationId}` | Get reservation details |
| DELETE | `/api/v1/reservations/{reservationId}` | Cancel a reservation |
| **Stock Movements (Audit)** |||
| GET | `/api/v1/movements` | List stock movements with filters |
| GET | `/api/v1/movements/{movementId}` | Get movement details |
| **Alerts** |||
| GET | `/api/v1/alerts/low-stock` | Get products below threshold |
| PUT | `/api/v1/products/{productId}/alert-threshold` | Set low-stock alert threshold |
| **Health** |||
| GET | `/health` | Health check endpoint |
| GET | `/ready` | Readiness probe |

---

## 2. Go Structs for Request and Response Payloads

```go
// file: internal/api/types/common.go
package types

import (
	"time"
)

// Pagination contains pagination parameters for list requests.
type Pagination struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginatedResponse wraps list responses with pagination metadata.
type PaginatedResponse[T any] struct {
	Data       []T              `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta contains pagination metadata in responses.
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// Timestamp embeds common timestamp fields.
type Timestamp struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AuditInfo contains audit trail information.
type AuditInfo struct {
	PerformedBy   string `json:"performed_by"`
	PerformedAt   time.Time `json:"performed_at"`
	CorrelationID string `json:"correlation_id,omitempty"`
	Source        string `json:"source"` // "api", "kafka", "system"
}
```

```go
// file: internal/api/types/products.go
package types

import (
	"github.com/google/uuid"
)

// VariantAttribute represents a product variant attribute (e.g., size, color).
type VariantAttribute struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Value string `json:"value" validate:"required,min=1,max=100"`
}

// Product represents a product in the inventory system.
type Product struct {
	ID                 uuid.UUID          `json:"id"`
	SKU                string             `json:"sku"`
	Name               string             `json:"name"`
	Description        string             `json:"description,omitempty"`
	Category           string             `json:"category,omitempty"`
	LowStockThreshold  int                `json:"low_stock_threshold"`
	IsActive           bool               `json:"is_active"`
	Variants           []ProductVariant   `json:"variants,omitempty"`
	Timestamp
}

// ProductVariant represents a specific variant of a product.
type ProductVariant struct {
	ID          uuid.UUID          `json:"id"`
	ProductID   uuid.UUID          `json:"product_id"`
	SKU         string             `json:"sku"`
	Attributes  []VariantAttribute `json:"attributes"`
	IsActive    bool               `json:"is_active"`
	Timestamp
}

// CreateProductRequest is the payload for creating a new product.
// @Description Request body for creating a new product
type CreateProductRequest struct {
	// SKU is the unique stock keeping unit identifier
	SKU string `json:"sku" validate:"required,min=3,max=50,alphanum"`
	
	// Name is the human-readable product name
	Name string `json:"name" validate:"required,min=1,max=255"`
	
	// Description provides additional product details
	Description string `json:"description,omitempty" validate:"max=2000"`
	
	// Category for product classification
	Category string `json:"category,omitempty" validate:"max=100"`
	
	// LowStockThreshold triggers alerts when stock falls below this level
	LowStockThreshold int `json:"low_stock_threshold" validate:"min=0"`
	
	// Variants are optional product variants to create with the product
	Variants []CreateVariantRequest `json:"variants,omitempty" validate:"dive"`
}

// UpdateProductRequest is the payload for updating a product.
type UpdateProductRequest struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description       *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	Category          *string `json:"category,omitempty" validate:"omitempty,max=100"`
	LowStockThreshold *int    `json:"low_stock_threshold,omitempty" validate:"omitempty,min=0"`
	IsActive          *bool   `json:"is_active,omitempty"`
}

// CreateVariantRequest is the payload for creating a product variant.
type CreateVariantRequest struct {
	SKU        string             `json:"sku" validate:"required,min=3,max=50"`
	Attributes []VariantAttribute `json:"attributes" validate:"required,min=1,dive"`
}

// ProductResponse is the API response for a product.
type ProductResponse struct {
	Product
	TotalStock      int `json:"total_stock"`
	AvailableStock  int `json:"available_stock"`
	ReservedStock   int `json:"reserved_stock"`
}

// ProductListFilters contains filters for listing products.
type ProductListFilters struct {
	Pagination
	SKU       string `json:"sku,omitempty"`
	Name      string `json:"name,omitempty"`
	Category  string `json:"category,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	LowStock  *bool  `json:"low_stock,omitempty"` // Filter for products below threshold
}
```

```go
// file: internal/api/types/warehouses.go
package types

import (
	"github.com/google/uuid"
)

// Warehouse represents a storage location.
type Warehouse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   Address   `json:"address"`
	IsActive  bool      `json:"is_active"`
	Priority  int       `json:"priority"` // Lower number = higher priority for allocation
	Timestamp
}

// Address represents a physical address.
type Address struct {
	Street     string `json:"street" validate:"required,max=255"`
	City       string `json:"city" validate:"required,max=100"`
	State      string `json:"state" validate:"required,max=100"`
	PostalCode string `json:"postal_code" validate:"required,max=20"`
	Country    string `json:"country" validate:"required,iso3166_1_alpha2"`
}

// CreateWarehouseRequest is the payload for creating a warehouse.
type CreateWarehouseRequest struct {
	Code     string  `json:"code" validate:"required,min=2,max=20,alphanum"`
	Name     string  `json:"name" validate:"required,min=1,max=255"`
	Address  Address `json:"address" validate:"required"`
	Priority int     `json:"priority" validate:"min=1,max=100"`
}

// UpdateWarehouseRequest is the payload for updating a warehouse.
type UpdateWarehouseRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Address  *Address `json:"address,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
	Priority *int     `json:"priority,omitempty" validate:"omitempty,min=1,max=100"`
}

// WarehouseResponse includes warehouse with stock summary.
type WarehouseResponse struct {
	Warehouse
	TotalProducts  int `json:"total_products"`
	TotalStock     int `json:"total_stock"`
	ReservedStock  int `json:"reserved_stock"`
}
```

```go
// file: internal/api/types/stock.go
package types

import (
	"time"

	"github.com/google/uuid"
)

// StockItem represents inventory level for a product/variant in a warehouse.
type StockItem struct {
	ID           uuid.UUID  `json:"id"`
	ProductID    uuid.UUID  `json:"product_id"`
	VariantID    *uuid.UUID `json:"variant_id,omitempty"`
	WarehouseID  uuid.UUID  `json:"warehouse_id"`
	Quantity     int        `json:"quantity"`      // Total physical stock
	Reserved     int        `json:"reserved"`      // Reserved for orders
	Available    int        `json:"available"`     // quantity - reserved
	Version      int        `json:"version"`       // Optimistic locking
	LastCountAt  *time.Time `json:"last_count_at,omitempty"`
	Timestamp
}

// AggregatedStock represents stock levels across all warehouses.
type AggregatedStock struct {
	ProductID       uuid.UUID             `json:"product_id"`
	VariantID       *uuid.UUID            `json:"variant_id,omitempty"`
	TotalQuantity   int                   `json:"total_quantity"`
	TotalReserved   int                   `json:"total_reserved"`
	TotalAvailable  int                   `json:"total_available"`
	ByWarehouse     []WarehouseStockLevel `json:"by_warehouse"`
	IsLowStock      bool                  `json:"is_low_stock"`
	Threshold       int                   `json:"threshold"`
}

// WarehouseStockLevel represents stock in a specific warehouse.
type WarehouseStockLevel struct {
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	WarehouseCode string    `json:"warehouse_code"`
	WarehouseName string    `json:"warehouse_name"`
	Quantity      int       `json:"quantity"`
	Reserved      int       `json:"reserved"`
	Available     int       `json:"available"`
}

// StockListFilters contains filters for listing stock.
type StockListFilters struct {
	Pagination
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	VariantID   *uuid.UUID `json:"variant_id,omitempty"`
	WarehouseID *uuid.UUID `json:"warehouse_id,omitempty"`
	LowStock    *bool      `json:"low_stock,omitempty"`
	MinQuantity *int       `json:"min_quantity,omitempty"`
	MaxQuantity *int       `json:"max_quantity,omitempty"`
}

// SetStockRequest is used to directly set stock level (admin operation).
type SetStockRequest struct {
	Quantity int       `json:"quantity" validate:"min=0"`
	Reason   string    `json:"reason" validate:"required,min=5,max=500"`
	AuditInfo
}
```

```go
// file: internal/api/types/stock_operations.go
package types

import (
	"github.com/google/uuid"
)

// StockReservationItem represents a single line item in a reservation request.
type StockReservationItem struct {
	ProductID   uuid.UUID  `json:"product_id" validate:"required"`
	VariantID   *uuid.UUID `json:"variant_id,omitempty"`
	Quantity    int        `json:"quantity" validate:"required,min=1"`
	WarehouseID *uuid.UUID `json:"warehouse_id,omitempty"` // Optional: system allocates if not specified
}

// ReserveStockRequest is the payload for reserving stock.
// @Description Reserve stock for an order. Supports multiple items and automatic warehouse allocation.
type ReserveStockRequest struct {
	// OrderID is the unique identifier of the order requesting stock
	OrderID string `json:"order_id" validate:"required,min=1,max=100"`
	
	// Items to reserve
	Items []StockReservationItem `json:"items" validate:"required,min=1,dive"`
	
	// ExpiresInMinutes sets reservation expiry (default: 30 minutes)
	ExpiresInMinutes int `json:"expires_in_minutes,omitempty" validate:"omitempty,min=1,max=1440"`
	
	// AllocationStrategy determines how stock is allocated across warehouses
	// Options: "single" (one warehouse), "distributed" (minimize shipping)
	AllocationStrategy string `json:"allocation_strategy,omitempty" validate:"omitempty,oneof=single distributed"`
	
	AuditInfo
}

// ReserveStockResponse contains the result of a reservation request.
type ReserveStockResponse struct {
	ReservationID string                    `json:"reservation_id"`
	OrderID       string                    `json:"order_id"`
	Status        string                    `json:"status"` // "confirmed", "partial", "failed"
	Items         []ReservationItemResult   `json:"items"`
	ExpiresAt     string                    `json:"expires_at"`
}

// ReservationItemResult contains the result for a single reservation item.
type ReservationItemResult struct {
	ProductID       uuid.UUID              `json:"product_id"`
	VariantID       *uuid.UUID             `json:"variant_id,omitempty"`
	RequestedQty    int                    `json:"requested_quantity"`
	ReservedQty     int                    `json:"reserved_quantity"`
	Status          string                 `json:"status"` // "reserved", "partial", "insufficient"
	Allocations     []WarehouseAllocation  `json:"allocations"`
	ShortageAmount  int                    `json:"shortage_amount,omitempty"`
}

// WarehouseAllocation represents stock allocated from a specific warehouse.
type WarehouseAllocation struct {
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	WarehouseCode string    `json:"warehouse_code"`
	Quantity      int       `json:"quantity"`
}

// ReleaseStockRequest is the payload for releasing reserved stock.
type ReleaseStockRequest struct {
	// ReservationID to release (mutually exclusive with OrderID)
	ReservationID string `json:"reservation_id,omitempty" validate:"required_without=OrderID"`
	
	// OrderID to release all reservations for an order
	OrderID string `json:"order_id,omitempty" validate:"required_without=ReservationID"`
	
	// Items to partially release (optional - releases all if not specified)
	Items []ReleaseItem `json:"items,omitempty" validate:"omitempty,dive"`
	
	// Reason for releasing the stock
	Reason string `json:"reason" validate:"required,min=5,max=500"`
	
	AuditInfo
}

// ReleaseItem specifies quantity to release for a specific product.
type ReleaseItem struct {
	ProductID uuid.UUID  `json:"product_id" validate:"required"`
	VariantID *uuid.UUID `json:"variant_id,omitempty"`
	Quantity  int        `json:"quantity" validate:"required,min=1"`
}

// ReleaseStockResponse contains the result of a release operation.
type ReleaseStockResponse struct {
	ReservationID    string `json:"reservation_id"`
	ReleasedItems    int    `json:"released_items"`
	ReleasedQuantity int    `json:"released_quantity"`
	Status           string `json:"status"` // "released", "partial", "not_found"
}

// FulfillStockRequest decrements stock when an order is shipped.
type FulfillStockRequest struct {
	// ReservationID of the reservation to fulfill
	ReservationID string `json:"reservation_id" validate:"required"`
	
	// OrderID for correlation
	OrderID string `json:"order_id" validate:"required"`
	
	// Items to fulfill (optional - fulfills all reserved items if not specified)
	Items []FulfillItem `json:"items,omitempty" validate:"omitempty,dive"`
	
	// ShipmentID for tracking
	ShipmentID string `json:"shipment_id,omitempty"`
	
	AuditInfo
}

// FulfillItem specifies quantity to fulfill for a specific product.
type FulfillItem struct {
	ProductID   uuid.UUID  `json:"product_id" validate:"required"`
	VariantID   *uuid.UUID `json:"variant_id,omitempty"`
	WarehouseID uuid.UUID  `json:"warehouse_id" validate:"required"`
	Quantity    int        `json:"quantity" validate:"required,min=1"`
}

// FulfillStockResponse contains the result of a fulfillment operation.
type FulfillStockResponse struct {
	ReservationID     string `json:"reservation_id"`
	OrderID           string `json:"order_id"`
	FulfilledItems    int    `json:"fulfilled_items"`
	FulfilledQuantity int    `json:"fulfilled_quantity"`
	Status            string `json:"status"` // "fulfilled", "partial"
}

// ReplenishStockRequest adds stock to a warehouse.
type ReplenishStockRequest struct {
	ProductID    uuid.UUID  `json:"product_id" validate:"required"`
	VariantID    *uuid.UUID `json:"variant_id,omitempty"`
	WarehouseID  uuid.UUID  `json:"warehouse_id" validate:"required"`
	Quantity     int        `json:"quantity" validate:"required,min=1"`
	
	// Reference to external system (e.g., purchase order)
	ReferenceType string `json:"reference_type,omitempty" validate:"omitempty,oneof=purchase_order return adjustment"`
	ReferenceID   string `json:"reference_id,omitempty"`
	
	// BatchInfo for tracking lot/batch numbers
	BatchNumber   string `json:"batch_number,omitempty"`
	ExpiryDate    string `json:"expiry_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	
	AuditInfo
}

// ReplenishStockResponse contains the result of replenishment.
type ReplenishStockResponse struct {
	StockItem
	PreviousQuantity int    `json:"previous_quantity"`
	AddedQuantity    int    `json:"added_quantity"`
	MovementID       string `json:"movement_id"`
}

// TransferStockRequest moves stock between warehouses.
type TransferStockRequest struct {
	ProductID         uuid.UUID  `json:"product_id" validate:"required"`
	VariantID         *uuid.UUID `json:"variant_id,omitempty"`
	SourceWarehouseID uuid.UUID  `json:"source_warehouse_id" validate:"required"`
	TargetWarehouseID uuid.UUID  `json:"target_warehouse_id" validate:"required,nefield=SourceWarehouseID"`
	Quantity          int        `json:"quantity" validate:"required,min=1"`
	
	// Reference for tracking
	TransferReference string `json:"transfer_reference,omitempty"`
	
	AuditInfo
}

// TransferStockResponse contains the result of a transfer.
type TransferStockResponse struct {
	TransferID          string    `json:"transfer_id"`
	ProductID           uuid.UUID `json:"product_id"`
	SourceWarehouseID   uuid.UUID `json:"source_warehouse_id"`
	TargetWarehouseID   uuid.UUID `json:"target_warehouse_id"`
	Quantity            int       `json:"quantity"`
	SourceNewQuantity   int       `json:"source_new_quantity"`
	TargetNewQuantity   int       `json:"target_new_quantity"`
	MovementIDs         []string  `json:"movement_ids"`
}
```

```go
// file: internal/api/types/reservations.go
package types

import (
	"time"

	"github.com/google/uuid"
)

// ReservationStatus represents the state of a reservation.
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "pending"
	ReservationStatusConfirmed ReservationStatus = "confirmed"
	ReservationStatusPartial   ReservationStatus = "partial"
	ReservationStatusFulfilled ReservationStatus = "fulfilled"
	ReservationStatusReleased  ReservationStatus = "released"
	ReservationStatusExpired   ReservationStatus = "expired"
)

// Reservation represents a stock reservation for an order.
type Reservation struct {
	ID        uuid.UUID         `json:"id"`
	OrderID   string            `json:"order_id"`
	Status    ReservationStatus `json:"status"`
	Items     []ReservationItem `json:"items"`
	ExpiresAt time.Time         `json:"expires_at"`
	Timestamp
	AuditInfo
}

// ReservationItem represents a single item in a reservation.
type ReservationItem struct {
	ID          uuid.UUID             `json:"id"`
	ProductID   uuid.UUID             `json:"product_id"`
	VariantID   *uuid.UUID            `json:"variant_id,omitempty"`
	Quantity    int                   `json:"quantity"`
	Allocations []WarehouseAllocation `json:"allocations"`
	Status      string                `json:"status"`
}

// ReservationListFilters contains filters for listing reservations.
type ReservationListFilters struct {
	Pagination
	OrderID     string             `json:"order_id,omitempty"`
	ProductID   *uuid.UUID         `json:"product_id,omitempty"`
	WarehouseID *uuid.UUID         `json:"warehouse_id,omitempty"`
	Status      *ReservationStatus `json:"status,omitempty"`
	ExpiringIn  *int               `json:"expiring_in_minutes,omitempty"` // Filter reservations expiring within N minutes
	CreatedFrom *time.Time         `json:"created_from,omitempty"`
	CreatedTo   *time.Time         `json:"created_to,omitempty"`
}
```

```go
// file: internal/api/types/movements.go
package types

import (
	"time"

	"github.com/google/uuid"
)

// MovementType represents the type of stock movement.
type MovementType string

const (
	MovementTypeReserve    MovementType = "reserve"
	MovementTypeRelease    MovementType = "release"
	MovementTypeFulfill    MovementType = "fulfill"
	MovementTypeReplenish  MovementType = "replenish"
	MovementTypeTransferIn MovementType = "transfer_in"
	MovementTypeTransferOut MovementType = "transfer_out"
	MovementTypeAdjustment MovementType = "adjustment"
	MovementTypeCount      MovementType = "count"
)

// StockMovement represents an auditable stock change event.
type StockMovement struct {
	ID              uuid.UUID    `json:"id"`
	ProductID       uuid.UUID    `json:"product_id"`
	VariantID       *uuid.UUID   `json:"variant_id,omitempty"`
	WarehouseID     uuid.UUID    `json:"warehouse_id"`
	MovementType    MovementType `json:"movement_type"`
	Quantity        int          `json:"quantity"`        // Positive for additions, negative for removals
	PreviousQty     int          `json:"previous_quantity"`
	NewQty          int          `json:"new_quantity"`
	PreviousReserved int         `json:"previous_reserved"`
	NewReserved     int          `json:"new_reserved"`
	
	// Reference fields for correlation
	ReservationID   *uuid.UUID `json:"reservation_id,omitempty"`
	OrderID         *string    `json:"order_id,omitempty"`
	TransferID      *uuid.UUID `json:"transfer_id,omitempty"`
	ReferenceType   *string    `json:"reference_type,omitempty"`
	ReferenceID     *string    `json:"reference_id,omitempty"`
	
	Reason          string     `json:"reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	AuditInfo
}

// MovementListFilters contains filters for listing movements.
type MovementListFilters struct {
	Pagination
	ProductID     *uuid.UUID    `json:"product_id,omitempty"`
	VariantID     *uuid.UUID    `json:"variant_id,omitempty"`
	WarehouseID   *uuid.UUID    `json:"warehouse_id,omitempty"`
	MovementType  *MovementType `json:"movement_type,omitempty"`
	OrderID       *string       `json:"order_id,omitempty"`
	ReservationID *uuid.UUID    `json:"reservation_id,omitempty"`
	PerformedBy   *string       `json:"performed_by,omitempty"`
	DateFrom      *time.Time    `json:"date_from,omitempty"`
	DateTo        *time.Time    `json:"date_to,omitempty"`
}

// MovementSummary provides aggregated movement statistics.
type MovementSummary struct {
	ProductID       uuid.UUID                  `json:"product_id"`
	WarehouseID     *uuid.UUID                 `json:"warehouse_id,omitempty"`
	Period          string                     `json:"period"` // "day", "week", "month"
	TotalIn         int                        `json:"total_in"`
	TotalOut        int                        `json:"total_out"`
	NetChange       int                        `json:"net_change"`
	ByType          map[MovementType]int       `json:"by_type"`
}
```

```go
// file: internal/api/types/alerts.go
package types

import (
	"time"

	"github.com/google/uuid"
)

// LowStockAlert represents a product that is below its threshold.
type LowStockAlert struct {
	ProductID         uuid.UUID  `json:"product_id"`
	ProductSKU        string     `json:"product_sku"`
	ProductName       string     `json:"product_name"`
	VariantID         *uuid.UUID `json:"variant_id,omitempty"`
	VariantSKU        *string    `json:"variant_sku,omitempty"`
	Threshold         int        `json:"threshold"`
	CurrentStock      int        `json:"current_stock"`
	AvailableStock    int        `json:"available_stock"`
	ShortageAmount    int        `json:"shortage_amount"`
	ByWarehouse       []WarehouseStockLevel `json:"by_warehouse"`
	LastReplenishedAt *time.Time `json:"last_replenished_at,omitempty"`
	AlertTriggeredAt  time.Time  `json:"alert_triggered_at"`
}

// LowStockAlertFilters contains filters for low stock alerts.
type LowStockAlertFilters struct {
	Pagination
	Category     string     `json:"category,omitempty"`
	WarehouseID  *uuid.UUID `json:"warehouse_id,omitempty"`
	MinShortage  *int       `json:"min_shortage,omitempty"`
	SortBy       string     `json:"sort_by,omitempty" validate:"omitempty,oneof=shortage_amount current_stock product_name"`
	SortOrder    string     `json:"sort_order,omitempty" validate:"omitempty,oneof=asc desc"`
}

// SetAlertThresholdRequest updates the low-stock threshold for a product.
type SetAlertThresholdRequest struct {
	Threshold int    `json:"threshold" validate:"min=0"`
	Reason    string `json:"reason,omitempty"`
	AuditInfo
}
```

```go
// file: internal/api/types/errors.go
package types

import (
	"time"
)

// ErrorCode represents application-specific error codes.
type ErrorCode string

const (
	// General errors
	ErrCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrCodeValidation          ErrorCode = "VALIDATION_ERROR"
	ErrCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrCodeConflict            ErrorCode = "CONFLICT"
	ErrCodeInternalError       ErrorCode = "INTERNAL_ERROR"
	ErrCodeServiceUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"
	
	// Domain-specific errors
	ErrCodeInsufficientStock   ErrorCode = "INSUFFICIENT_STOCK"
	ErrCodeReservationExpired  ErrorCode = "RESERVATION_EXPIRED"
	ErrCodeReservationNotFound ErrorCode = "RESERVATION_NOT_FOUND"
	ErrCodeProductInactive     ErrorCode = "PRODUCT_INACTIVE"
	ErrCodeWarehouseInactive   ErrorCode = "WAREHOUSE_INACTIVE"
	ErrCodeDuplicateSKU        ErrorCode = "DUPLICATE_SKU"
	ErrCodeStockNegative       ErrorCode = "STOCK_CANNOT_BE_NEGATIVE"
	ErrCodeOptimisticLock      ErrorCode = "CONCURRENT_MODIFICATION"
	ErrCodeTransferSameWarehouse ErrorCode = "TRANSFER_SAME_WAREHOUSE"
)

// ErrorResponse is the standard error response format.
// @Description Standard error response format for all API errors
type ErrorResponse struct {
	// Error contains the primary error information
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information.
type ErrorDetail struct {
	// Code is the application-specific error