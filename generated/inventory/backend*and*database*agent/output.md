# Backend & Database Agent Output

# Inventory Microservice Implementation

I'll implement a complete inventory microservice following Clean Architecture principles with all the specified entities, operations, and API contracts.

## Domain Layer

### Entities

```go
// file: internal/domain/entity/product.go
package entity

import (
	"errors"
	"time"
)

// ProductVariant represents a specific variant of a product (e.g., size, color)
type ProductVariant struct {
	Size  string
	Color string
}

// Product represents a product in the inventory system
type Product struct {
	ID          string
	SKU         string
	Name        string
	Description string
	Variant     ProductVariant
	Category    string
	MinStock    int // Threshold for low-stock alerts
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Validation errors
var (
	ErrProductIDRequired   = errors.New("product ID is required")
	ErrProductSKURequired  = errors.New("product SKU is required")
	ErrProductNameRequired = errors.New("product name is required")
	ErrMinStockNegative    = errors.New("minimum stock cannot be negative")
	ErrProductDeleted      = errors.New("product has been deleted")
)

// NewProduct creates a new Product with validation
func NewProduct(id, sku, name, description, category string, variant ProductVariant, minStock int) (*Product, error) {
	if id == "" {
		return nil, ErrProductIDRequired
	}
	if sku == "" {
		return nil, ErrProductSKURequired
	}
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if minStock < 0 {
		return nil, ErrMinStockNegative
	}

	now := time.Now().UTC()
	return &Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		Description: description,
		Variant:     variant,
		Category:    category,
		MinStock:    minStock,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update modifies product details
func (p *Product) Update(name, description, category string, variant ProductVariant, minStock int) error {
	if p.DeletedAt != nil {
		return ErrProductDeleted
	}
	if name == "" {
		return ErrProductNameRequired
	}
	if minStock < 0 {
		return ErrMinStockNegative
	}

	p.Name = name
	p.Description = description
	p.Category = category
	p.Variant = variant
	p.MinStock = minStock
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// SoftDelete marks the product as deleted
func (p *Product) SoftDelete() error {
	if p.DeletedAt != nil {
		return ErrProductDeleted
	}
	now := time.Now().UTC()
	p.DeletedAt = &now
	p.IsActive = false
	p.UpdatedAt = now
	return nil
}

// IsDeleted returns true if the product has been soft deleted
func (p *Product) IsDeleted() bool {
	return p.DeletedAt != nil
}
```

```go
// file: internal/domain/entity/warehouse.go
package entity

import (
	"errors"
	"time"
)

// WarehouseAddress represents the physical address of a warehouse
type WarehouseAddress struct {
	Street     string
	City       string
	State      string
	Country    string
	PostalCode string
}

// Warehouse represents a storage location for inventory
type Warehouse struct {
	ID        string
	Code      string
	Name      string
	Address   WarehouseAddress
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Warehouse validation errors
var (
	ErrWarehouseIDRequired   = errors.New("warehouse ID is required")
	ErrWarehouseCodeRequired = errors.New("warehouse code is required")
	ErrWarehouseNameRequired = errors.New("warehouse name is required")
	ErrWarehouseDeleted      = errors.New("warehouse has been deleted")
)

// NewWarehouse creates a new Warehouse with validation
func NewWarehouse(id, code, name string, address WarehouseAddress) (*Warehouse, error) {
	if id == "" {
		return nil, ErrWarehouseIDRequired
	}
	if code == "" {
		return nil, ErrWarehouseCodeRequired
	}
	if name == "" {
		return nil, ErrWarehouseNameRequired
	}

	now := time.Now().UTC()
	return &Warehouse{
		ID:        id,
		Code:      code,
		Name:      name,
		Address:   address,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update modifies warehouse details
func (w *Warehouse) Update(name string, address WarehouseAddress) error {
	if w.DeletedAt != nil {
		return ErrWarehouseDeleted
	}
	if name == "" {
		return ErrWarehouseNameRequired
	}

	w.Name = name
	w.Address = address
	w.UpdatedAt = time.Now().UTC()
	return nil
}

// SoftDelete marks the warehouse as deleted
func (w *Warehouse) SoftDelete() error {
	if w.DeletedAt != nil {
		return ErrWarehouseDeleted
	}
	now := time.Now().UTC()
	w.DeletedAt = &now
	w.IsActive = false
	w.UpdatedAt = now
	return nil
}

// IsDeleted returns true if the warehouse has been soft deleted
func (w *Warehouse) IsDeleted() bool {
	return w.DeletedAt != nil
}
```

```go
// file: internal/domain/entity/stock_item.go
package entity

import (
	"errors"
	"time"
)

// StockItem represents the stock level of a product in a specific warehouse
type StockItem struct {
	ID              string
	ProductID       string
	WarehouseID     string
	QuantityOnHand  int // Physical stock available
	QuantityReserved int // Stock reserved for pending orders
	ReorderPoint    int // When to trigger replenishment
	ReorderQuantity int // How much to reorder
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// StockItem validation errors
var (
	ErrStockItemIDRequired      = errors.New("stock item ID is required")
	ErrStockItemProductRequired = errors.New("product ID is required")
	ErrStockItemWarehouseRequired = errors.New("warehouse ID is required")
	ErrQuantityNegative         = errors.New("quantity cannot be negative")
	ErrInsufficientStock        = errors.New("insufficient stock available")
	ErrInsufficientReserved     = errors.New("insufficient reserved stock")
	ErrReorderPointNegative     = errors.New("reorder point cannot be negative")
	ErrReorderQuantityNegative  = errors.New("reorder quantity cannot be negative")
)

// NewStockItem creates a new StockItem with validation
func NewStockItem(id, productID, warehouseID string, reorderPoint, reorderQuantity int) (*StockItem, error) {
	if id == "" {
		return nil, ErrStockItemIDRequired
	}
	if productID == "" {
		return nil, ErrStockItemProductRequired
	}
	if warehouseID == "" {
		return nil, ErrStockItemWarehouseRequired
	}
	if reorderPoint < 0 {
		return nil, ErrReorderPointNegative
	}
	if reorderQuantity < 0 {
		return nil, ErrReorderQuantityNegative
	}

	now := time.Now().UTC()
	return &StockItem{
		ID:              id,
		ProductID:       productID,
		WarehouseID:     warehouseID,
		QuantityOnHand:  0,
		QuantityReserved: 0,
		ReorderPoint:    reorderPoint,
		ReorderQuantity: reorderQuantity,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// AvailableQuantity returns the quantity available for reservation
func (s *StockItem) AvailableQuantity() int {
	return s.QuantityOnHand - s.QuantityReserved
}

// Reserve attempts to reserve a quantity of stock
func (s *StockItem) Reserve(quantity int) error {
	if quantity < 0 {
		return ErrQuantityNegative
	}
	if s.AvailableQuantity() < quantity {
		return ErrInsufficientStock
	}

	s.QuantityReserved += quantity
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// ReleaseReservation releases previously reserved stock
func (s *StockItem) ReleaseReservation(quantity int) error {
	if quantity < 0 {
		return ErrQuantityNegative
	}
	if s.QuantityReserved < quantity {
		return ErrInsufficientReserved
	}

	s.QuantityReserved -= quantity
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// Fulfill decrements both reserved and on-hand quantities (order shipped)
func (s *StockItem) Fulfill(quantity int) error {
	if quantity < 0 {
		return ErrQuantityNegative
	}
	if s.QuantityReserved < quantity {
		return ErrInsufficientReserved
	}
	if s.QuantityOnHand < quantity {
		return ErrInsufficientStock
	}

	s.QuantityReserved -= quantity
	s.QuantityOnHand -= quantity
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// Replenish adds stock to the on-hand quantity
func (s *StockItem) Replenish(quantity int) error {
	if quantity < 0 {
		return ErrQuantityNegative
	}

	s.QuantityOnHand += quantity
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// NeedsReorder returns true if stock is at or below reorder point
func (s *StockItem) NeedsReorder() bool {
	return s.AvailableQuantity() <= s.ReorderPoint
}

// IsLowStock returns true if available quantity is below or equal to reorder point
func (s *StockItem) IsLowStock() bool {
	return s.AvailableQuantity() <= s.ReorderPoint
}
```

```go
// file: internal/domain/entity/stock_movement.go
package entity

import (
	"errors"
	"time"
)

// MovementType represents the type of stock movement
type MovementType string

const (
	MovementTypeReplenishment MovementType = "REPLENISHMENT"
	MovementTypeReservation   MovementType = "RESERVATION"
	MovementTypeRelease       MovementType = "RELEASE"
	MovementTypeFulfillment   MovementType = "FULFILLMENT"
	MovementTypeAdjustment    MovementType = "ADJUSTMENT"
	MovementTypeTransfer      MovementType = "TRANSFER"
)

// StockMovement represents an audit record of stock changes
type StockMovement struct {
	ID              string
	StockItemID     string
	MovementType    MovementType
	Quantity        int    // Positive for additions, negative for reductions
	ReferenceID     string // Order ID, reservation ID, etc.
	ReferenceType   string // "ORDER", "RESERVATION", "MANUAL", etc.
	PreviousOnHand  int
	NewOnHand       int
	PreviousReserved int
	NewReserved     int
	Reason          string
	CreatedBy       string
	CreatedAt       time.Time
}

// StockMovement validation errors
var (
	ErrMovementIDRequired        = errors.New("movement ID is required")
	ErrMovementStockItemRequired = errors.New("stock item ID is required")
	ErrMovementTypeInvalid       = errors.New("invalid movement type")
	ErrMovementQuantityZero      = errors.New("movement quantity cannot be zero")
)

// NewStockMovement creates a new StockMovement with validation
func NewStockMovement(
	id, stockItemID string,
	movementType MovementType,
	quantity int,
	referenceID, referenceType string,
	previousOnHand, newOnHand int,
	previousReserved, newReserved int,
	reason, createdBy string,
) (*StockMovement, error) {
	if id == "" {
		return nil, ErrMovementIDRequired
	}
	if stockItemID == "" {
		return nil, ErrMovementStockItemRequired
	}
	if !isValidMovementType(movementType) {
		return nil, ErrMovementTypeInvalid
	}
	if quantity == 0 {
		return nil, ErrMovementQuantityZero
	}

	return &StockMovement{
		ID:              id,
		StockItemID:     stockItemID,
		MovementType:    movementType,
		Quantity:        quantity,
		ReferenceID:     referenceID,
		ReferenceType:   referenceType,
		PreviousOnHand:  previousOnHand,
		NewOnHand:       newOnHand,
		PreviousReserved: previousReserved,
		NewReserved:     newReserved,
		Reason:          reason,
		CreatedBy:       createdBy,
		CreatedAt:       time.Now().UTC(),
	}, nil
}

func isValidMovementType(mt MovementType) bool {
	switch mt {
	case MovementTypeReplenishment, MovementTypeReservation, MovementTypeRelease,
		MovementTypeFulfillment, MovementTypeAdjustment, MovementTypeTransfer:
		return true
	}
	return false
}
```

```go
// file: internal/domain/entity/reservation.go
package entity

import (
	"errors"
	"time"
)

// ReservationStatus represents the current state of a reservation
type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusReleased  ReservationStatus = "RELEASED"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

// ReservationItem represents a single item in a reservation
type ReservationItem struct {
	StockItemID string
	ProductID   string
	WarehouseID string
	Quantity    int
}

// Reservation represents stock reserved for an order
type Reservation struct {
	ID          string
	OrderID     string
	Items       []ReservationItem
	Status      ReservationStatus
	ExpiresAt   time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ReleasedAt  *time.Time
	FulfilledAt *time.Time
}

// Reservation validation errors
var (
	ErrReservationIDRequired     = errors.New("reservation ID is required")
	ErrReservationOrderRequired  = errors.New("order ID is required")
	ErrReservationItemsRequired  = errors.New("at least one reservation item is required")
	ErrReservationItemQuantity   = errors.New("reservation item quantity must be positive")
	ErrReservationNotPending     = errors.New("reservation is not in pending status")
	ErrReservationNotConfirmed   = errors.New("reservation is not in confirmed status")
	ErrReservationAlreadyReleased = errors.New("reservation has already been released")
	ErrReservationAlreadyFulfilled = errors.New("reservation has already been fulfilled")
	ErrReservationExpired        = errors.New("reservation has expired")
)

// NewReservation creates a new Reservation with validation
func NewReservation(id, orderID string, items []ReservationItem, expiresAt time.Time) (*Reservation, error) {
	if id == "" {
		return nil, ErrReservationIDRequired
	}
	if orderID == "" {
		return nil, ErrReservationOrderRequired
	}
	if len(items) == 0 {
		return nil, ErrReservationItemsRequired
	}
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrReservationItemQuantity
		}
	}

	now := time.Now().UTC()
	return &Reservation{
		ID:        id,
		OrderID:   orderID,
		Items:     items,
		Status:    ReservationStatusPending,
		ExpiresAt: expiresAt,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Confirm transitions the reservation to confirmed status
func (r *Reservation) Confirm() error {
	if r.Status != ReservationStatusPending {
		return ErrReservationNotPending
	}
	if time.Now().UTC().After(r.ExpiresAt) {
		return ErrReservationExpired
	}

	r.Status = ReservationStatusConfirmed
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// Release releases the reserved stock back to available
func (r *Reservation) Release() error {
	if r.Status == ReservationStatusReleased {
		return ErrReservationAlreadyReleased
	}
	if r.Status == ReservationStatusFulfilled {
		return ErrReservationAlreadyFulfilled
	}

	now := time.Now().UTC()
	r.Status = ReservationStatusReleased
	r.ReleasedAt = &now
	r.UpdatedAt = now
	return nil
}

// Fulfill marks the reservation as fulfilled (order shipped)
func (r *Reservation) Fulfill() error {
	if r.Status == ReservationStatusFulfilled {
		return ErrReservationAlreadyFulfilled
	}
	if r.Status == ReservationStatusReleased {
		return ErrReservationAlreadyReleased
	}
	if r.Status != ReservationStatusConfirmed && r.Status != ReservationStatusPending {
		return ErrReservationNotConfirmed
	}

	now := time.Now().UTC()
	r.Status = ReservationStatusFulfilled
	r.FulfilledAt = &now
	r.UpdatedAt = now
	return nil
}

// Expire marks the reservation as expired
func (r *Reservation) Expire() error {
	if r.Status != ReservationStatusPending && r.Status != ReservationStatusConfirmed {
		return ErrReservationNotPending
	}

	r.Status = ReservationStatusExpired
	r.UpdatedAt = time.Now().UTC()
	return nil
}

// IsExpired checks if the reservation has expired
func (r *Reservation) IsExpired() bool {
	return time.Now().UTC().After(r.ExpiresAt) && r.Status != ReservationStatusFulfilled && r.Status != ReservationStatusReleased
}

// TotalQuantity returns the total quantity across all items
func (r *Reservation) TotalQuantity() int {
	total := 0
	for _, item := range r.Items {
		total += item.Quantity
	}
	return total
}
```

```go
// file: internal/domain/entity/low_stock_alert.go
package entity

import (
	"errors"
	"time"
)

// AlertStatus represents the current state of a low stock alert
type AlertStatus string

const (
	AlertStatusActive      AlertStatus = "ACTIVE"
	AlertStatusAcknowledged AlertStatus = "ACKNOWLEDGED"
	AlertStatusResolved    AlertStatus = "RESOLVED"
)

// LowStockAlert represents an alert for low stock levels
type LowStockAlert struct {
	ID              string
	StockItemID     string
	ProductID       string
	WarehouseID     string
	CurrentQuantity int
	ReorderPoint    int
	Status          AlertStatus
	AcknowledgedBy  *string
	AcknowledgedAt  *time.Time
	ResolvedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// LowStockAlert validation errors
var (
	ErrAlertIDRequired        = errors.New("alert ID is required")
	ErrAlertStockItemRequired = errors.New("stock item ID is required")
	ErrAlertProductRequired   = errors.New("product ID is required")
	ErrAlertWarehouseRequired = errors.New("warehouse ID is required")
	ErrAlertAlreadyResolved   = errors.New("alert has already been resolved")
)

// NewLowStockAlert creates a new LowStockAlert
func NewLowStockAlert(id, stockItemID, productID, warehouseID string, currentQuantity, reorderPoint int) (*LowStockAlert, error) {
	if id == "" {
		return nil, ErrAlertIDRequired
	}
	if stockItemID == "" {
		return nil, ErrAlertStockItemRequired
	}
	if productID == "" {
		return nil, ErrAlertProductRequired
	}
	if warehouseID == "" {
		return nil, ErrAlertWarehouseRequired
	}

	now := time.Now().UTC()
	return &LowStockAlert{
		ID:              id,
		StockItemID:     stockItemID,
		ProductID:       productID,
		WarehouseID:     warehouseID,
		CurrentQuantity: currentQuantity,
		ReorderPoint:    reorderPoint,
		Status:          AlertStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Acknowledge marks the alert as acknowledged
func (a *LowStockAlert) Acknowledge(userID string) error {
	if a.Status == AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	now := time.Now().UTC()
	a.Status = AlertStatusAcknowledged
	a.AcknowledgedBy = &userID
	a.AcknowledgedAt = &now
	a.UpdatedAt = now
	return nil
}

// Resolve marks the alert as resolved
func (a *LowStockAlert) Resolve() error {
	if a.Status == AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	now := time.Now().UTC()
	a.Status = AlertStatusResolved
	a.ResolvedAt = &now
	a.UpdatedAt = now
	return nil
}
```

### Repository Interfaces

```go
// file: internal/domain/repository/product_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// ProductFilter defines filtering options for product queries
type ProductFilter struct {
	SKU      *string
	Name     *string
	Category *string
	IsActive *bool
	Limit    int
	Offset   int
}

// ProductRepository defines the interface for product persistence
type ProductRepository interface {
	// Create persists a new product
	Create(ctx context.Context, product *entity.Product) error

	// GetByID retrieves a product by its ID
	GetByID(ctx context.Context, id string) (*entity.Product, error)

	// GetBySKU retrieves a product by its SKU
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)

	// List retrieves products with optional filtering
	List(ctx context.Context, filter ProductFilter) ([]*entity.Product, int, error)

	// Update persists changes to an existing product
	Update(ctx context.Context, product *entity.Product) error

	// Delete soft deletes a product
	Delete(ctx context.Context, id string) error

	// ExistsBySKU checks if a product with the given SKU exists
	ExistsBySKU(ctx context.Context, sku string) (bool, error)
}
```

```go
// file: internal/domain/repository/warehouse_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// WarehouseFilter defines filtering options for warehouse queries
type WarehouseFilter struct {
	Code     *string
	Name     *string
	IsActive *bool
	Limit    int
	Offset   int
}

// WarehouseRepository defines the interface for warehouse persistence
type WarehouseRepository interface {
	// Create persists a new warehouse
	Create(ctx context.Context, warehouse *entity.Warehouse) error

	// GetByID retrieves a warehouse by its ID
	GetByID(ctx context.Context, id string) (*entity.Warehouse, error)

	// GetByCode retrieves a warehouse by its code
	GetByCode(ctx context.Context, code string) (*entity.Warehouse, error)

	// List retrieves warehouses with optional filtering
	List(ctx context.Context, filter WarehouseFilter) ([]*entity.Warehouse, int, error)

	// Update persists changes to an existing warehouse
	Update(ctx context.Context, warehouse *entity.Warehouse) error

	// Delete soft deletes a warehouse
	Delete(ctx context.Context, id string) error

	// ExistsByCode checks if a warehouse with the given code exists
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
```

```go
// file: internal/domain/repository/stock_item_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// StockItemFilter defines filtering options for stock item queries
type StockItemFilter struct {
	ProductID   *string
	WarehouseID *string
	LowStock    *bool // Filter items at or below reorder point
	Limit       int
	Offset      int
}

// AggregatedStock represents total stock for a product across warehouses
type AggregatedStock struct {
	ProductID        string
	TotalOnHand      int
	TotalReserved    int
	TotalAvailable   int
	WarehouseCount   int
	WarehouseDetails []WarehouseStockDetail
}

// WarehouseStockDetail represents stock in a specific warehouse
type WarehouseStockDetail struct {
	WarehouseID      string
	WarehouseName    string
	QuantityOnHand   int
	QuantityReserved int
	Available        int
}

// StockItemRepository defines the interface for stock item persistence
type StockItemRepository interface {
	// Create persists a new stock item
	Create(ctx context.Context, stockItem *entity.StockItem) error

	// GetByID retrieves a stock item by its ID
	GetByID(ctx context.Context, id string) (*entity.StockItem, error)

	// GetByProductAndWarehouse retrieves a stock item by product and warehouse
	GetByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*entity.StockItem, error)

	// List retrieves stock items with optional filtering
	List(ctx context.Context, filter StockItemFilter) ([]*entity.StockItem, int, error)

	// Update persists changes to an existing stock item
	Update(ctx context.Context, stockItem *entity.StockItem) error

	// UpdateWithLock updates a stock item with optimistic locking
	UpdateWithLock(ctx context.Context, stockItem *entity.StockItem, expectedVersion int) error

	// GetAggregatedStock retrieves total stock for a product across all warehouses
	GetAggregatedStock(ctx context.Context, productID string) (*AggregatedStock, error)

	// GetLowStockItems retrieves all stock items at or below reorder point
	GetLowStockItems(ctx context.Context) ([]*entity.StockItem, error)

	// ExistsByProductAndWarehouse checks if a stock item exists for the given product and warehouse
	ExistsByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (bool, error)
}
```

```go
// file: internal/domain/repository/stock_movement_repository.go
package repository

import (
	"context"
	"time"

	"github.com/inventory-service/internal/domain/entity"
)

// StockMovementFilter defines filtering options for stock movement queries
type StockMovementFilter struct {
	StockItemID   *string
	MovementType  *entity.MovementType
	ReferenceID   *string
	ReferenceType *string
	StartDate     *time.Time
	EndDate       *time.Time
	Limit         int
	Offset        int
}

// StockMovementRepository defines the interface for stock movement persistence
type StockMovementRepository interface {
	// Create persists a new stock movement record
	Create(ctx context.Context, movement *entity.StockMovement) error

	// GetByID retrieves a stock movement by its ID
	GetByID(ctx context.Context, id string) (*entity.StockMovement, error)

	// List retrieves stock movements with optional filtering
	List(ctx context.Context, filter StockMovementFilter) ([]*entity.StockMovement, int, error)

	// GetByStockItem retrieves all movements for a specific stock item
	GetByStockItem(ctx context.Context, stockItemID string, limit, offset int) ([]*entity.StockMovement, int, error)

	// GetByReference retrieves movements by reference (e.g., order ID)
	GetByReference(ctx context.Context, referenceID, referenceType string) ([]*entity.StockMovement, error)
}
```

```go
// file: internal/domain/repository/reservation_repository.go
package repository

import (
	"context"
	"time"

	"github.com/inventory-service/internal/domain/entity"
)

// ReservationFilter defines filtering options for reservation queries
type ReservationFilter struct {
	OrderID   *string
	Status    *entity.ReservationStatus
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

// ReservationRepository defines the interface for reservation persistence
type ReservationRepository interface {
	// Create persists a new reservation
	Create(ctx context.Context, reservation *entity.Reservation) error

	// GetByID retrieves a reservation by its ID
	GetByID(ctx context.Context, id string) (*entity.Reservation, error)

	// GetByOrderID retrieves reservations for a specific order
	GetByOrderID(ctx context.Context, orderID string) ([]*entity.Reservation, error)

	// List retrieves reservations with optional filtering
	List(ctx context.Context, filter ReservationFilter) ([]*entity.Reservation, int, error)

	// Update persists changes to an existing reservation
	Update(ctx context.Context, reservation *entity.Reservation) error

	// GetExpiredReservations retrieves all reservations that have expired
	GetExpiredReservations(ctx context.Context) ([]*entity.Reservation, error)
}
```

```go
// file: internal/domain/repository/low_stock_alert_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// LowStockAlertFilter defines filtering options for alert queries
type LowStockAlertFilter struct {
	ProductID   *string
	WarehouseID *string
	Status      *entity.AlertStatus
	Limit       int
	Offset      int
}

// LowStockAlertRepository defines the interface for low stock alert persistence
type LowStockAlertRepository interface {
	// Create persists a new low stock alert
	Create(ctx context.Context, alert *entity.LowStockAlert) error

	// GetByID retrieves an alert by its ID
	GetByID(ctx context.Context, id string) (*entity.Low