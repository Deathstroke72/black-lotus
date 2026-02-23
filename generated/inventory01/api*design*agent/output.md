# API Design Agent Output

# Inventory Microservice API Design

## 1. Complete API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| **Products** |
| GET | `/api/v1/products` | List all products with pagination |
| POST | `/api/v1/products` | Create a new product |
| GET | `/api/v1/products/{productId}` | Get product details |
| PUT | `/api/v1/products/{productId}` | Update product information |
| DELETE | `/api/v1/products/{productId}` | Soft delete a product |
| GET | `/api/v1/products/{productId}/variants` | List all variants for a product |
| POST | `/api/v1/products/{productId}/variants` | Create a product variant |
| **Warehouses** |
| GET | `/api/v1/warehouses` | List all warehouses |
| POST | `/api/v1/warehouses` | Create a new warehouse |
| GET | `/api/v1/warehouses/{warehouseId}` | Get warehouse details |
| PUT | `/api/v1/warehouses/{warehouseId}` | Update warehouse information |
| DELETE | `/api/v1/warehouses/{warehouseId}` | Deactivate a warehouse |
| **Stock Items** |
| GET | `/api/v1/stock` | Query stock levels with filters |
| GET | `/api/v1/stock/{productId}` | Get aggregated stock for a product |
| GET | `/api/v1/stock/{productId}/warehouses` | Get stock breakdown by warehouse |
| GET | `/api/v1/warehouses/{warehouseId}/stock` | Get all stock in a warehouse |
| PUT | `/api/v1/warehouses/{warehouseId}/stock/{productId}` | Set stock level (admin override) |
| **Stock Operations** |
| POST | `/api/v1/stock/reserve` | Reserve stock for an order |
| POST | `/api/v1/stock/release` | Release reserved stock |
| POST | `/api/v1/stock/fulfill` | Decrement stock on fulfillment |
| POST | `/api/v1/stock/replenish` | Replenish stock |
| POST | `/api/v1/stock/transfer` | Transfer stock between warehouses |
| **Stock Movements (Audit)** |
| GET | `/api/v1/movements` | Query stock movements with filters |
| GET | `/api/v1/movements/{movementId}` | Get movement details |
| GET | `/api/v1/products/{productId}/movements` | Get movements for a product |
| **Reservations** |
| GET | `/api/v1/reservations` | List reservations with filters |
| GET | `/api/v1/reservations/{reservationId}` | Get reservation details |
| GET | `/api/v1/reservations/order/{orderId}` | Get reservations for an order |
| DELETE | `/api/v1/reservations/{reservationId}` | Cancel a specific reservation |
| **Alerts** |
| GET | `/api/v1/alerts/low-stock` | Get current low-stock alerts |
| PUT | `/api/v1/products/{productId}/thresholds` | Set low-stock thresholds |
| **Health** |
| GET | `/health` | Service health check |
| GET | `/ready` | Readiness probe |

---

## 2. Go Structs for Request/Response Payloads

```go
// file: internal/api/types/common.go
package types

import (
	"time"
)

// Pagination represents pagination parameters for list endpoints
type Pagination struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginatedResponse wraps list responses with pagination metadata
type PaginatedResponse[T any] struct {
	Data       []T              `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta contains pagination metadata in responses
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// Timestamp is a wrapper for consistent time formatting
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).UTC().Format(time.RFC3339) + `"`), nil
}

// AuditInfo contains common audit fields
type AuditInfo struct {
	CreatedAt Timestamp  `json:"created_at"`
	UpdatedAt Timestamp  `json:"updated_at"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
}
```

```go
// file: internal/api/types/errors.go
package types

// ErrorResponse is the standardized error response format
// @Description Standard error response returned by all endpoints
type ErrorResponse struct {
	// Error contains the main error information
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	// Code is a machine-readable error code
	Code string `json:"code"`
	// Message is a human-readable error description
	Message string `json:"message"`
	// Details contains additional error context
	Details map[string]interface{} `json:"details,omitempty"`
	// ValidationErrors contains field-level validation errors
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	// TraceID for request tracing
	TraceID string `json:"trace_id,omitempty"`
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

// Standard error codes
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInsufficientStock = "INSUFFICIENT_STOCK"
	ErrCodeReservationExpired = "RESERVATION_EXPIRED"
	ErrCodeInvalidState     = "INVALID_STATE"
	ErrCodeInternalError    = "INTERNAL_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
	ErrCodeForbidden        = "FORBIDDEN"
)
```

```go
// file: internal/api/types/product.go
package types

import (
	"github.com/google/uuid"
)

// VariantAttribute represents a product variant attribute (size, color, etc.)
type VariantAttribute struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Value string `json:"value" validate:"required,min=1,max=100"`
}

// Product represents a product in the inventory system
type Product struct {
	ID          uuid.UUID          `json:"id"`
	SKU         string             `json:"sku"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Category    string             `json:"category,omitempty"`
	Variants    []VariantAttribute `json:"variants,omitempty"`
	IsActive    bool               `json:"is_active"`
	LowStockThreshold int          `json:"low_stock_threshold"`
	AuditInfo
}

// CreateProductRequest represents the request body for creating a product
// @Description Request payload for creating a new product
type CreateProductRequest struct {
	// SKU is the unique stock keeping unit identifier
	SKU string `json:"sku" validate:"required,min=3,max=50,alphanumunicode"`
	// Name is the product display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Description is an optional product description
	Description string `json:"description,omitempty" validate:"max=2000"`
	// Category is the product category
	Category string `json:"category,omitempty" validate:"max=100"`
	// LowStockThreshold triggers alerts when stock falls below this level
	LowStockThreshold int `json:"low_stock_threshold" validate:"min=0"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description       *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	Category          *string `json:"category,omitempty" validate:"omitempty,max=100"`
	IsActive          *bool   `json:"is_active,omitempty"`
	LowStockThreshold *int    `json:"low_stock_threshold,omitempty" validate:"omitempty,min=0"`
}

// CreateVariantRequest represents the request body for creating a product variant
type CreateVariantRequest struct {
	// SKU is the unique identifier for this variant
	SKU string `json:"sku" validate:"required,min=3,max=50"`
	// Attributes define the variant properties (e.g., size: "L", color: "red")
	Attributes []VariantAttribute `json:"attributes" validate:"required,min=1,dive"`
}

// ProductResponse is the response wrapper for product operations
type ProductResponse struct {
	Product Product `json:"product"`
}

// ProductListResponse is the response for listing products
type ProductListResponse = PaginatedResponse[Product]

// ProductListParams contains query parameters for listing products
type ProductListParams struct {
	Pagination
	Category  string `json:"category,omitempty"`
	SKUPrefix string `json:"sku_prefix,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	Search    string `json:"search,omitempty"`
}
```

```go
// file: internal/api/types/warehouse.go
package types

import (
	"github.com/google/uuid"
)

// Address represents a physical address
type Address struct {
	Street     string `json:"street" validate:"required,max=255"`
	City       string `json:"city" validate:"required,max=100"`
	State      string `json:"state" validate:"max=100"`
	PostalCode string `json:"postal_code" validate:"required,max=20"`
	Country    string `json:"country" validate:"required,iso3166_1_alpha2"`
}

// Warehouse represents a storage location
type Warehouse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   Address   `json:"address"`
	IsActive  bool      `json:"is_active"`
	Priority  int       `json:"priority"`
	AuditInfo
}

// CreateWarehouseRequest represents the request body for creating a warehouse
type CreateWarehouseRequest struct {
	// Code is a unique short identifier for the warehouse
	Code string `json:"code" validate:"required,min=2,max=20,alphanum"`
	// Name is the warehouse display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Address is the physical location
	Address Address `json:"address" validate:"required"`
	// Priority determines fulfillment order (lower = higher priority)
	Priority int `json:"priority" validate:"min=0,max=1000"`
}

// UpdateWarehouseRequest represents the request body for updating a warehouse
type UpdateWarehouseRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Address  *Address `json:"address,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
	Priority *int     `json:"priority,omitempty" validate:"omitempty,min=0,max=1000"`
}

// WarehouseResponse is the response wrapper for warehouse operations
type WarehouseResponse struct {
	Warehouse Warehouse `json:"warehouse"`
}

// WarehouseListResponse is the response for listing warehouses
type WarehouseListResponse = PaginatedResponse[Warehouse]
```

```go
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
```

```go
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
```

```go
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
```

```go
// file: internal/api/types/reservation.go
package types

import (
	"github.com/google/uuid"
)

// ReservationStatus represents the status of a reservation
type ReservationStatus string

const (
	ReservationStatusActive    ReservationStatus = "ACTIVE"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

// Reservation represents a stock reservation
type Reservation struct {
	ID            uuid.UUID         `json:"id"`
	OrderID       string            `json:"order_id"`
	ProductID     uuid.UUID         `json:"product_id"`
	ProductSKU    string            `json:"product_sku"`
	WarehouseID   uuid.UUID         `json:"warehouse_id"`
	WarehouseCode string            `json:"warehouse_code"`
	Quantity      int               `json:"quantity"`
	Status        ReservationStatus `json:"status"`
	ExpiresAt     Timestamp         `json:"expires_at"`
	FulfilledAt   *Timestamp        `json:"fulfilled_at,omitempty"`
	ReleasedAt    *Timestamp        `json:"released_at,omitempty"`
	AuditInfo
}

// ReservationQueryParams contains query parameters for reservation queries
type ReservationQueryParams struct {
	Pagination
	OrderID     *string            `json:"order_id,omitempty"`
	ProductID   *uuid.UUID         `json:"product_id,omitempty"`
	WarehouseID *uuid.UUID         `json:"warehouse_id,omitempty"`
	Status      *ReservationStatus `json:"status,omitempty"`
	ExpiringSoon bool              `json:"expiring_soon,omitempty"` // Within 5 minutes
}

// ReservationResponse wraps single reservation responses
type ReservationResponse struct {
	Reservation Reservation `json:"reservation"`
}

// ReservationListResponse is the response for listing reservations
type ReservationListResponse = PaginatedResponse[Reservation]
```

```go
// file: internal/api/types/alerts.go
package types

import (
	"github.com/google/uuid"
)

// AlertSeverity represents the severity of a low-stock alert
type AlertSeverity string

const (
	AlertSeverityWarning  AlertSeverity = "WARNING"  // Below threshold
	AlertSeverityCritical AlertSeverity = "CRITICAL" // Below 25% of threshold
	AlertSeverityOutOfStock AlertSeverity = "OUT_OF_STOCK"
)

// LowStockAlert represents a low-stock alert
type LowStockAlert struct {
	ID              uuid.UUID     `json:"id"`
	ProductID       uuid.UUID     `json:"product_id"`
	ProductSKU      string        `json:"product_sku"`
	ProductName     string        `json:"product_name"`
	WarehouseID     *uuid.UUID    `json:"warehouse_id,omitempty"` // Null for aggregate alerts
	WarehouseCode   *string       `json:"warehouse_code,omitempty"`
	CurrentQuantity int           `json:"current_quantity"`
	Threshold       int           `json:"threshold"`
	Severity        AlertSeverity `json:"severity"`
	CreatedAt       Timestamp     `json:"created_at"`
	AcknowledgedAt  *Timestamp    `json:"acknowledged_at,omitempty"`
	AcknowledgedBy  *string       `json:"acknowledged_by,omitempty"`
}

// LowStockAlertParams contains query parameters for alert queries
type LowStockAlertParams struct {
	Pagination
	Severity        *AlertSeverity `json:"severity,omitempty"`
	WarehouseID     *uuid.UUID     `json:"warehouse_id,omitempty"`
	Category        *string        `json:"category,omitempty"`
	Unacknowledged  bool           `json:"unacknowledged,omitempty"`
}

// LowStockAlertListResponse is the response for listing alerts
type LowStockAlertListResponse = PaginatedResponse[LowStockAlert]

// SetThresholdsRequest sets low-stock thresholds for a product
type SetThresholdsRequest struct {
	// GlobalThreshold applies to aggregated stock
	GlobalThreshold int `json:"global_threshold" validate:"min=0"`
	// WarehouseThresholds are warehouse-specific overrides
	WarehouseThresholds map[uuid.UUID]int `json:"warehouse_thresholds,omitempty"`
}

// ThresholdsResponse shows current thresholds for a product
type ThresholdsResponse struct {
	ProductID           uuid.UUID         `json:"product_id"`
	GlobalThreshold     int               `json:"global_threshold"`
	WarehouseThresholds map[uuid.UUID]int `json:"warehouse_thresholds"`
}
```

---

## 3. Router Setup Code

```go
// file: internal/api/router/router.go
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"inventory-service/internal/api/handlers"
	apimiddleware "inventory-service/internal/api/middleware"
)

// Config holds router configuration
type Config struct {
	Environment     string
	RequestTimeout  time.Duration
	AllowedOrigins  []string
}

// Handlers aggregates all handler dependencies
type Handlers struct {
	Product     *handlers.ProductHandler
	Warehouse   *handlers.WarehouseHandler
	Stock       *handlers.StockHandler
	StockOps    *handlers.StockOperationsHandler
	Movement    *handlers.MovementHandler
	Reservation *handlers.ReservationHandler
	Alert       *handlers.AlertHandler
	Health      *handlers.HealthHandler
}

// NewRouter creates and configures the main router
func NewRouter(cfg Config, h Handlers) http.Handler {
	r := chi.NewRouter()

	// Global middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(apimiddleware.Logger())
	r.Use(apimiddleware.Recoverer())
	r.Use(middleware.Timeout(cfg.RequestTimeout))

	// CORS configuration
	r.Use(cors.