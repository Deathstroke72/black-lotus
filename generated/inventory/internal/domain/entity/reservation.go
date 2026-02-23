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