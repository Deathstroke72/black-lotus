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