# Testing & Security Agent Output

# Inventory Microservice - Tests and Security Implementation

I'll provide a comprehensive implementation of tests and security middleware for the inventory microservice.

## Project Structure

```
inventory/
├── internal/
│   ├── api/types/
│   ├── domain/
│   ├── repository/
│   ├── service/
│   ├── middleware/
│   └── testutil/
├── tests/
│   ├── unit/
│   ├── integration/
│   └── security/
├── Makefile
└── go.mod
```

## 1. Mock Generation and Test Utilities

```go
// file: internal/testutil/mocks.go
package testutil

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"inventory/internal/domain"
	"inventory/internal/repository"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{}
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context, filter repository.ProductFilter) ([]*domain.Product, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Product), args.Int(1), args.Error(2)
}

// MockStockRepository is a mock implementation of StockRepository
type MockStockRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockStockRepository() *MockStockRepository {
	return &MockStockRepository{}
}

func (m *MockStockRepository) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID, warehouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StockItem), args.Error(1)
}

func (m *MockStockRepository) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockItem), args.Error(1)
}

func (m *MockStockRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, warehouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockItem), args.Error(1)
}

func (m *MockStockRepository) Create(ctx context.Context, item *domain.StockItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockStockRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity, reserved int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity, reserved)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicReserve(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicRelease(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicFulfill(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicReplenish(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) GetLowStock(ctx context.Context, threshold int) ([]*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockItem), args.Error(1)
}

// MockReservationRepository is a mock implementation of ReservationRepository
type MockReservationRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockReservationRepository() *MockReservationRepository {
	return &MockReservationRepository{}
}

func (m *MockReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, reservation)
	return args.Error(0)
}

func (m *MockReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReservationStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReservationRepository) List(ctx context.Context, filter repository.ReservationFilter) ([]*domain.Reservation, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Reservation), args.Int(1), args.Error(2)
}

func (m *MockReservationRepository) GetExpired(ctx context.Context, before time.Time) ([]*domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, before)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

// MockStockMovementRepository is a mock implementation of StockMovementRepository
type MockStockMovementRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockStockMovementRepository() *MockStockMovementRepository {
	return &MockStockMovementRepository{}
}

func (m *MockStockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, movement)
	return args.Error(0)
}

func (m *MockStockMovementRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StockMovement), args.Error(1)
}

func (m *MockStockMovementRepository) List(ctx context.Context, filter repository.MovementFilter) ([]*domain.StockMovement, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.StockMovement), args.Int(1), args.Error(2)
}

func (m *MockStockMovementRepository) GetByProduct(ctx context.Context, productID uuid.UUID, limit int) ([]*domain.StockMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockMovement), args.Error(1)
}

// MockWarehouseRepository is a mock implementation of WarehouseRepository
type MockWarehouseRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockWarehouseRepository() *MockWarehouseRepository {
	return &MockWarehouseRepository{}
}

func (m *MockWarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) GetByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWarehouseRepository) List(ctx context.Context, activeOnly bool) ([]*domain.Warehouse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Warehouse), args.Error(1)
}

// MockCacheRepository is a mock implementation of CacheRepository
type MockCacheRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockCacheRepository() *MockCacheRepository {
	return &MockCacheRepository{}
}

func (m *MockCacheRepository) GetStockLevel(ctx context.Context, productID uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID)
	return args.Int(0), args.Error(1)
}

func (m *MockCacheRepository) SetStockLevel(ctx context.Context, productID uuid.UUID, quantity int, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, productID, quantity, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) InvalidateStockLevel(ctx context.Context, productID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, productID)
	return args.Error(0)
}

// MockEventPublisher is a mock implementation of EventPublisher
type MockEventPublisher struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

func (m *MockEventPublisher) PublishStockReserved(ctx context.Context, event domain.StockReservedEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishStockReleased(ctx context.Context, event domain.StockReleasedEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishLowStockAlert(ctx context.Context, event domain.LowStockAlertEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishStockFulfilled(ctx context.Context, event domain.StockFulfilledEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockTransactionManager is a mock implementation of TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

func NewMockTransactionManager() *MockTransactionManager {
	return &MockTransactionManager{}
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	// Execute the function to test transaction behavior
	if args.Bool(1) {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return args.Error(0)
}
```

```go
// file: internal/testutil/fixtures.go
package testutil

import (
	"time"

	"github.com/google/uuid"

	"inventory/internal/domain"
)

// TestFixtures provides common test data
type TestFixtures struct{}

func NewTestFixtures() *TestFixtures {
	return &TestFixtures{}
}

func (f *TestFixtures) Product(opts ...func(*domain.Product)) *domain.Product {
	p := &domain.Product{
		ID:                uuid.New(),
		SKU:               "TEST-SKU-001",
		Name:              "Test Product",
		Description:       "A test product for unit testing",
		Category:          "Electronics",
		BasePrice:         99.99,
		LowStockThreshold: 10,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (f *TestFixtures) Warehouse(opts ...func(*domain.Warehouse)) *domain.Warehouse {
	w := &domain.Warehouse{
		ID:        uuid.New(),
		Code:      "WH-001",
		Name:      "Main Warehouse",
		Address:   domain.Address{Street: "123 Main St", City: "Test City", Country: "US"},
		IsActive:  true,
		Priority:  1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (f *TestFixtures) StockItem(productID, warehouseID uuid.UUID, opts ...func(*domain.StockItem)) *domain.StockItem {
	s := &domain.StockItem{
		ID:               uuid.New(),
		ProductID:        productID,
		WarehouseID:      warehouseID,
		Quantity:         100,
		ReservedQuantity: 0,
		ReorderPoint:     10,
		ReorderQuantity:  50,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (f *TestFixtures) Reservation(productID, warehouseID, orderID uuid.UUID, opts ...func(*domain.Reservation)) *domain.Reservation {
	r := &domain.Reservation{
		ID:          uuid.New(),
		ProductID:   productID,
		WarehouseID: warehouseID,
		OrderID:     orderID,
		Quantity:    5,
		Status:      domain.ReservationStatusPending,
		ExpiresAt:   time.Now().Add(30 * time.Minute),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (f *TestFixtures) StockMovement(productID, warehouseID uuid.UUID, movementType domain.MovementType, opts ...func(*domain.StockMovement)) *domain.StockMovement {
	m := &domain.StockMovement{
		ID:               uuid.New(),
		ProductID:        productID,
		WarehouseID:      warehouseID,
		MovementType:     movementType,
		Quantity:         10,
		PreviousQuantity: 100,
		NewQuantity:      110,
		Reference:        "TEST-REF-001",
		CreatedAt:        time.Now(),
		CreatedBy:        uuid.New(),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// WithQuantity sets stock item quantity
func WithQuantity(qty int) func(*domain.StockItem) {
	return func(s *domain.StockItem) {
		s.Quantity = qty
	}
}

// WithReserved sets stock item reserved quantity
func WithReserved(qty int) func(*domain.StockItem) {
	return func(s *domain.StockItem) {
		s.ReservedQuantity = qty
	}
}

// WithSKU sets product SKU
func WithSKU(sku string) func(*domain.Product) {
	return func(p *domain.Product) {
		p.SKU = sku
	}
}

// WithThreshold sets product low stock threshold
func WithThreshold(threshold int) func(*domain.Product) {
	return func(p *domain.Product) {
		p.LowStockThreshold = threshold
	}
}

// WithInactive sets entity as inactive
func WithProductInactive() func(*domain.Product) {
	return func(p *domain.Product) {
		p.IsActive = false
	}
}

func WithWarehouseInactive() func(*domain.Warehouse) {
	return func(w *domain.Warehouse) {
		w.IsActive = false
	}
}

// WithReservationStatus sets reservation status
func WithReservationStatus(status domain.ReservationStatus) func(*domain.Reservation) {
	return func(r *domain.Reservation) {
		r.Status = status
	}
}

// WithExpiredReservation sets reservation as expired
func WithExpiredReservation() func(*domain.Reservation) {
	return func(r *domain.Reservation) {
		r.ExpiresAt = time.Now().Add(-1 * time.Hour)
	}
}
```

## 2. Domain Types

```go
// file: internal/domain/entities.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                uuid.UUID  `json:"id"`
	SKU               string     `json:"sku"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	Category          string     `json:"category"`
	BasePrice         float64    `json:"base_price"`
	LowStockThreshold int        `json:"low_stock_threshold"`
	IsActive          bool       `json:"is_active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type ProductVariant struct {
	ID                uuid.UUID         `json:"id"`
	ProductID         uuid.UUID         `json:"product_id"`
	SKU               string            `json:"sku"`
	Name              string            `json:"name"`
	Attributes        map[string]string `json:"attributes"`
	PriceModifier     float64           `json:"price_modifier"`
	LowStockThreshold *int              `json:"low_stock_threshold,omitempty"`
	IsActive          bool              `json:"is_active"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty"`
}

type Warehouse struct {
	ID        uuid.UUID  `json:"id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	Address   Address    `json:"address"`
	IsActive  bool       `json:"is_active"`
	Priority  int        `json:"priority"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type StockItem struct {
	ID               uuid.UUID  `json:"id"`
	ProductID        uuid.UUID  `json:"product_id"`
	VariantID        *uuid.UUID `json:"variant_id,omitempty"`
	WarehouseID      uuid.UUID  `json:"warehouse_id"`
	Quantity         int        `json:"quantity"`
	ReservedQuantity int        `json:"reserved_quantity"`
	ReorderPoint     int        `json:"reorder_point"`
	ReorderQuantity  int        `json:"reorder_quantity"`
	LastCountedAt    *time.Time `json:"last_counted_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// AvailableQuantity returns the quantity available for reservation
func (s *StockItem) AvailableQuantity() int {
	return s.Quantity - s.ReservedQuantity
}

// IsLowStock checks if stock is below reorder point
func (s *StockItem) IsLowStock() bool {
	return s.AvailableQuantity() <= s.ReorderPoint
}

type MovementType string

const (
	MovementTypeReservation  MovementType = "RESERVATION"
	MovementTypeRelease      MovementType = "RELEASE"
	MovementTypeFulfillment  MovementType = "FULFILLMENT"
	MovementTypeReplenishment MovementType = "REPLENISHMENT"
	MovementTypeAdjustment   MovementType = "ADJUSTMENT"
	MovementTypeTransferIn   MovementType = "TRANSFER_IN"
	MovementTypeTransferOut  MovementType = "TRANSFER_OUT"
)

type StockMovement struct {
	ID               uuid.UUID    `json:"id"`
	ProductID        uuid.UUID    `json:"product_id"`
	VariantID        *uuid.UUID   `json:"variant_id,omitempty"`
	WarehouseID      uuid.UUID    `json:"warehouse_id"`
	MovementType     MovementType `json:"movement_type"`
	Quantity         int          `json:"quantity"`
	PreviousQuantity int          `json:"previous_quantity"`
	NewQuantity      int          `json:"new_quantity"`
	Reference        string       `json:"reference"`
	Notes            string       `json:"notes,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	CreatedBy        uuid.UUID    `json:"created_by"`
}

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusConfirmed ReservationStatus = "CONFIRMED"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusCancelled ReservationStatus = "CANCELLED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

type Reservation struct {
	ID          uuid.UUID         `json:"id"`
	ProductID   uuid.UUID         `json:"product_id"`
	VariantID   *uuid.UUID        `json:"variant_id,omitempty"`
	WarehouseID uuid.UUID         `json:"warehouse_id"`
	OrderID     uuid.UUID         `json:"order_id"`
	Quantity    int               `json:"quantity"`
	Status      ReservationStatus `json:"status"`
	ExpiresAt   time.Time         `json:"expires_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// IsExpired checks if reservation has expired
func (r *Reservation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// CanBeCancelled checks if reservation can be cancelled
func (r *Reservation) CanBeCancelled() bool {
	return r.Status == ReservationStatusPending || r.Status == ReservationStatusConfirmed
}
```

```go
// file: internal/domain/events.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type StockReservedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
}

type StockReleasedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	Reason        string    `json:"reason"`
	Timestamp     time.Time `json:"timestamp"`
}

type StockFulfilledEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
}

type LowStockAlertEvent struct {
	ProductID         uuid.UUID `json:"product_id"`
	WarehouseID       uuid.UUID `json:"warehouse_id"`
	CurrentQuantity   int       `json:"current_quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	ReorderPoint      int       `json:"reorder_point"`
	Timestamp         time.Time `json:"timestamp"`
}
```

```go
// file: internal/domain/errors.go
package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound              = errors.New("resource not found")
	ErrInsufficientStock     = errors.New("insufficient stock available")
	ErrInvalidQuantity       = errors.New("invalid quantity")
	ErrProductInactive       = errors.New("product is inactive")
	E