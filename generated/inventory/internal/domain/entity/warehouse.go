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