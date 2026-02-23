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