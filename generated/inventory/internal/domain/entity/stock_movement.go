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