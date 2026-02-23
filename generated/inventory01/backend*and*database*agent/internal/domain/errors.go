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