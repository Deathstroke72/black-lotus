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