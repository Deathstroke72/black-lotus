# Backend & Database Agent Output

# Inventory Microservice - Complete Backend Implementation

I'll provide a comprehensive implementation of the inventory microservice with all requested components.

## 1. Database Schema and Migrations

```sql
-- file: migrations/000001_initial_schema.up.sql

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Products table
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sku VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    base_price DECIMAL(12, 2),
    low_stock_threshold INTEGER NOT NULL DEFAULT 10,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_active ON products(is_active) WHERE deleted_at IS NULL;

-- Product variants table (size, color, etc.)
CREATE TABLE product_variants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    attributes JSONB NOT NULL DEFAULT '{}',
    price_modifier DECIMAL(12, 2) DEFAULT 0,
    low_stock_threshold INTEGER,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT unique_product_variant UNIQUE (product_id, sku)
);

CREATE INDEX idx_product_variants_product_id ON product_variants(product_id);
CREATE INDEX idx_product_variants_sku ON product_variants(sku);
CREATE INDEX idx_product_variants_attributes ON product_variants USING GIN(attributes);

-- Warehouses table
CREATE TABLE warehouses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    address JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT true,
    priority INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_warehouses_code ON warehouses(code);
CREATE INDEX idx_warehouses_is_active ON warehouses(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_warehouses_priority ON warehouses(priority DESC);

-- Stock items table
CREATE TABLE stock_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    reorder_point INTEGER NOT NULL DEFAULT 10,
    reorder_quantity INTEGER NOT NULL DEFAULT 100,
    last_counted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    CONSTRAINT unique_stock_item UNIQUE (product_id, variant_id, warehouse_id),
    CONSTRAINT non_negative_quantity CHECK (quantity >= 0),
    CONSTRAINT non_negative_reserved CHECK (reserved_quantity >= 0),
    CONSTRAINT reserved_not_exceeds_quantity CHECK (reserved_quantity <= quantity)
);

CREATE INDEX idx_stock_items_product_id ON stock_items(product_id);
CREATE INDEX idx_stock_items_variant_id ON stock_items(variant_id);
CREATE INDEX idx_stock_items_warehouse_id ON stock_items(warehouse_id);
CREATE INDEX idx_stock_items_low_stock ON stock_items((quantity - reserved_quantity)) 
    WHERE (quantity - reserved_quantity) <= reorder_point;

-- Stock movement types enum
CREATE TYPE stock_movement_type AS ENUM (
    'RESERVATION',
    'RELEASE',
    'FULFILLMENT',
    'REPLENISHMENT',
    'ADJUSTMENT',
    'TRANSFER_OUT',
    'TRANSFER_IN',
    'RETURN',
    'DAMAGE',
    'EXPIRED'
);

-- Stock movements table (audit trail)
CREATE TABLE stock_movements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stock_item_id UUID NOT NULL REFERENCES stock_items(id),
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    movement_type stock_movement_type NOT NULL,
    quantity INTEGER NOT NULL,
    quantity_before INTEGER NOT NULL,
    quantity_after INTEGER NOT NULL,
    reserved_before INTEGER NOT NULL,
    reserved_after INTEGER NOT NULL,
    reference_type VARCHAR(50),
    reference_id UUID,
    reason TEXT,
    metadata JSONB DEFAULT '{}',
    performed_by VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_stock_movements_stock_item_id ON stock_movements(stock_item_id);
CREATE INDEX idx_stock_movements_product_id ON stock_movements(product_id);
CREATE INDEX idx_stock_movements_warehouse_id ON stock_movements(warehouse_id);
CREATE INDEX idx_stock_movements_reference ON stock_movements(reference_type, reference_id);
CREATE INDEX idx_stock_movements_created_at ON stock_movements(created_at DESC);
CREATE INDEX idx_stock_movements_type ON stock_movements(movement_type);

-- Reservation status enum
CREATE TYPE reservation_status AS ENUM (
    'PENDING',
    'CONFIRMED',
    'PARTIALLY_FULFILLED',
    'FULFILLED',
    'CANCELLED',
    'EXPIRED'
);

-- Reservations table
CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    stock_item_id UUID NOT NULL REFERENCES stock_items(id),
    quantity INTEGER NOT NULL,
    fulfilled_quantity INTEGER NOT NULL DEFAULT 0,
    status reservation_status NOT NULL DEFAULT 'PENDING',
    expires_at TIMESTAMPTZ NOT NULL,
    confirmed_at TIMESTAMPTZ,
    fulfilled_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT fulfilled_not_exceeds_quantity CHECK (fulfilled_quantity <= quantity)
);

CREATE INDEX idx_reservations_order_id ON reservations(order_id);
CREATE INDEX idx_reservations_product_id ON reservations(product_id);
CREATE INDEX idx_reservations_warehouse_id ON reservations(warehouse_id);
CREATE INDEX idx_reservations_stock_item_id ON reservations(stock_item_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_expires_at ON reservations(expires_at) WHERE status IN ('PENDING', 'CONFIRMED');

-- Low stock alerts table
CREATE TABLE low_stock_alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id),
    variant_id UUID REFERENCES product_variants(id),
    warehouse_id UUID REFERENCES warehouses(id),
    current_quantity INTEGER NOT NULL,
    threshold INTEGER NOT NULL,
    is_acknowledged BOOLEAN NOT NULL DEFAULT false,
    acknowledged_by VARCHAR(100),
    acknowledged_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX idx_low_stock_alerts_product_id ON low_stock_alerts(product_id);
CREATE INDEX idx_low_stock_alerts_unresolved ON low_stock_alerts(created_at DESC) 
    WHERE resolved_at IS NULL;

-- Outbox table for reliable event publishing
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT
);

CREATE INDEX idx_outbox_events_unpublished ON outbox_events(created_at) 
    WHERE published_at IS NULL;
CREATE INDEX idx_outbox_events_aggregate ON outbox_events(aggregate_type, aggregate_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_product_variants_updated_at BEFORE UPDATE ON product_variants
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_warehouses_updated_at BEFORE UPDATE ON warehouses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_stock_items_updated_at BEFORE UPDATE ON stock_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

```sql
-- file: migrations/000001_initial_schema.down.sql

DROP TRIGGER IF EXISTS update_reservations_updated_at ON reservations;
DROP TRIGGER IF EXISTS update_stock_items_updated_at ON stock_items;
DROP TRIGGER IF EXISTS update_warehouses_updated_at ON warehouses;
DROP TRIGGER IF EXISTS update_product_variants_updated_at ON product_variants;
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS low_stock_alerts;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS warehouses;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;

DROP TYPE IF EXISTS reservation_status;
DROP TYPE IF EXISTS stock_movement_type;

DROP EXTENSION IF EXISTS "uuid-ossp";
```

## 2. Domain Models

```go
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
```

## 3. Custom Errors

```go
// file: internal/domain/errors.go
package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrNotFound              = errors.New("resource not found")
	ErrConflict              = errors.New("resource conflict")
	ErrInsufficientStock     = errors.New("insufficient stock available")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrReservationExpired    = errors.New("reservation has expired")
	ErrReservationNotPending = errors.New("reservation is not in pending status")
	ErrOptimisticLock        = errors.New("optimistic lock conflict, resource was modified")
	ErrWarehouseInactive     = errors.New("warehouse is inactive")
	ErrProductInactive       = errors.New("product is inactive")
	ErrInvalidOperation      = errors.New("invalid operation")
)

// NotFoundError represents a resource not found error
type NotFoundError struct {
	Resource string
	ID       uuid.UUID
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with id %s not found", e.Resource, e.ID)
}

func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// InsufficientStockError represents an insufficient stock error
type InsufficientStockError struct {
	ProductID   uuid.UUID
	VariantID   *uuid.UUID
	WarehouseID uuid.UUID
	Requested   int
	Available   int
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf(
		"insufficient stock for product %s in warehouse %s: requested %d, available %d",
		e.ProductID, e.WarehouseID, e.Requested, e.Available,
	)
}

func (e *InsufficientStockError) Is(target error) bool {
	return target == ErrInsufficientStock
}

// ConflictError represents a resource conflict error
type ConflictError struct {
	Resource string
	Field    string
	Value    string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("%s with %s '%s' already exists", e.Resource, e.Field, e.Value)
}

func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
```

## 4. Repository Interfaces

```go
// file: internal/repository/interfaces.go
package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/inventory-service/internal/domain"
)

// ProductFilter defines filters for querying products
type ProductFilter struct {
	SKU      *string
	Category *string
	IsActive *bool
	Page     int
	PageSize int
}

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetBySKU(ctx context.Context, sku string) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ProductFilter) ([]*domain.Product, int, error)
	GetVariants(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariant, error)
	CreateVariant(ctx context.Context, variant *domain.ProductVariant) error
	GetVariantByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error)
	UpdateVariant(ctx context.Context, variant *domain.ProductVariant) error
	SoftDeleteVariant(ctx context.Context, id uuid.UUID) error
}

// WarehouseFilter defines filters for querying warehouses
type WarehouseFilter struct {
	Code     *string
	IsActive *bool
	Page     int
	PageSize int
}

// WarehouseRepository defines the interface for warehouse data access
type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *domain.Warehouse) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*domain.Warehouse, error)
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter WarehouseFilter) ([]*domain.Warehouse, int, error)
	GetActiveByPriority(ctx context.Context) ([]*domain.Warehouse, error)
}

// StockFilter defines filters for querying stock items
type StockFilter struct {
	ProductID   *uuid.UUID
	VariantID   *uuid.UUID
	WarehouseID *uuid.UUID
	LowStock    *bool
	Page        int
	PageSize    int
}

// StockRepository defines the interface for stock item data access
type StockRepository interface {
	Create(ctx context.Context, item *domain.StockItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.StockItem, error)
	GetByProductWarehouse(ctx context.Context, productID, warehouseID uuid.UUID, variantID *uuid.UUID) (*domain.StockItem, error)
	Update(ctx context.Context, item *domain.StockItem) error
	UpdateWithVersion(ctx context.Context, item *domain.StockItem, expectedVersion int) error
	List(ctx context.Context, filter StockFilter) ([]*domain.StockItem, int, error)
	GetByProduct(ctx context.Context, productID uuid.UUID, variantID *uuid.UUID) ([]*domain.StockItem, error)
	GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.StockItem, error)
	GetLowStockItems(ctx context.Context) ([]*domain.StockItem, error)
	
	// Atomic operations with row-level locking
	ReserveStock(ctx context.Context, id uuid.UUID, quantity int) (*domain.StockItem, error)
	ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) (*domain.StockItem, error)
	FulfillStock(ctx context.Context, id uuid.UUID, quantity int) (*domain.StockItem, error)
	ReplenishStock(ctx context.Context, id uuid.UUID, quantity int) (*domain.StockItem, error)
	AdjustStock(ctx context.Context, id uuid.UUID, newQuantity int) (*domain.StockItem, error)
}

// MovementFilter defines filters for querying stock movements
type MovementFilter struct {
	ProductID     *uuid.UUID
	WarehouseID   *uuid.UUID
	MovementType  *domain.StockMovementType
	ReferenceType *string
	ReferenceID   *uuid.UUID
	StartDate     *time.Time
	EndDate       *time.Time
	Page          int
	PageSize      int
}

// StockMovementRepository defines the interface for stock movement data access
type StockMovementRepository interface {
	Create(ctx context.Context, movement *domain.StockMovement) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error)
	List(ctx context.Context, filter MovementFilter) ([]*domain.StockMovement, int, error)
	GetByProduct(ctx context.Context, productID uuid.UUID, limit int) ([]*domain.StockMovement, error)
	GetByReference(ctx context.Context, refType string, refID uuid.UUID) ([]*domain.StockMovement, error)
}

// ReservationFilter defines filters for querying reservations
type ReservationFilter struct {
	OrderID     *uuid.UUID
	ProductID   *uuid.UUID
	WarehouseID *uuid.UUID
	Status      *domain.ReservationStatus
	ExpiresBefore *time.Time
	Page        int
	PageSize    int
}

// ReservationRepository defines the interface for reservation data access
type ReservationRepository interface {
	Create(ctx context.Context, reservation *domain.Reservation) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error)
	Update(ctx context.Context, reservation *domain.Reservation) error
	List(ctx context.Context, filter ReservationFilter) ([]*domain.Reservation, int, error)
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Reservation, error)
	GetActiveByStockItem(ctx context.Context, stockItemID uuid.UUID) ([]*domain.Reservation, error)
	GetExpired(ctx context.Context) ([]*domain.Reservation, error)
	CancelExpired(ctx context.Context) (int, error)
}

// LowStockAlertRepository defines the interface for low stock alert data access
type LowStockAlertRepository interface {
	Create(ctx context.Context, alert *domain.LowStockAlert) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.LowStockAlert, error)
	GetUnresolved(ctx context.Context) ([]*domain.LowStockAlert, error)
	GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.LowStockAlert, error)
	Acknowledge(ctx context.Context