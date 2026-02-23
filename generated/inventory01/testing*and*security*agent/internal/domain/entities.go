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