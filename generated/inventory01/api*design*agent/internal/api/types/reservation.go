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