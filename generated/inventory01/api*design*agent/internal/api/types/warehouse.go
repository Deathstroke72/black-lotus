// file: internal/api/types/warehouse.go
package types

import (
	"github.com/google/uuid"
)

// Address represents a physical address
type Address struct {
	Street     string `json:"street" validate:"required,max=255"`
	City       string `json:"city" validate:"required,max=100"`
	State      string `json:"state" validate:"max=100"`
	PostalCode string `json:"postal_code" validate:"required,max=20"`
	Country    string `json:"country" validate:"required,iso3166_1_alpha2"`
}

// Warehouse represents a storage location
type Warehouse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   Address   `json:"address"`
	IsActive  bool      `json:"is_active"`
	Priority  int       `json:"priority"`
	AuditInfo
}

// CreateWarehouseRequest represents the request body for creating a warehouse
type CreateWarehouseRequest struct {
	// Code is a unique short identifier for the warehouse
	Code string `json:"code" validate:"required,min=2,max=20,alphanum"`
	// Name is the warehouse display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Address is the physical location
	Address Address `json:"address" validate:"required"`
	// Priority determines fulfillment order (lower = higher priority)
	Priority int `json:"priority" validate:"min=0,max=1000"`
}

// UpdateWarehouseRequest represents the request body for updating a warehouse
type UpdateWarehouseRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Address  *Address `json:"address,omitempty"`
	IsActive *bool    `json:"is_active,omitempty"`
	Priority *int     `json:"priority,omitempty" validate:"omitempty,min=0,max=1000"`
}

// WarehouseResponse is the response wrapper for warehouse operations
type WarehouseResponse struct {
	Warehouse Warehouse `json:"warehouse"`
}

// WarehouseListResponse is the response for listing warehouses
type WarehouseListResponse = PaginatedResponse[Warehouse]