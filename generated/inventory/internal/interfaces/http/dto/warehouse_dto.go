// file: internal/interfaces/http/dto/warehouse_dto.go
package dto

import "time"

// WarehouseAddress represents a warehouse physical address.
type WarehouseAddress struct {
	// Street is the street address
	Street string `json:"street" validate:"required,max=255"`
	// City is the city name
	City string `json:"city" validate:"required,max=100"`
	// State is the state or province
	State string `json:"state,omitempty" validate:"max=100"`
	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code" validate:"required,max=20"`
	// Country is the ISO 3166-1 alpha-2 country code
	Country string `json:"country" validate:"required,len=2"`
}

// CreateWarehouseRequest represents the request body for creating a warehouse.
// @Description Request payload for creating a new warehouse
type CreateWarehouseRequest struct {
	// Name is the warehouse display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Code is a unique warehouse code (e.g., "WH-NYC-01")
	Code string `json:"code" validate:"required,min=1,max=50"`
	// Address is the warehouse physical address
	Address WarehouseAddress `json:"address" validate:"required"`
	// IsActive indicates if the warehouse is operational
	IsActive bool `json:"is_active"`
	// Priority is used for stock allocation (lower = higher priority)
	Priority int `json:"priority" validate:"min=0"`
}

// UpdateWarehouseRequest represents the request body for updating a warehouse.
// @Description Request payload for updating an existing warehouse
type UpdateWarehouseRequest struct {
	// Name is the warehouse display name
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	// Address is the warehouse physical address
	Address *WarehouseAddress `json:"address,omitempty"`
	// IsActive indicates if the warehouse is operational
	IsActive *bool `json:"is_active,omitempty"`
	// Priority is used for stock allocation
	Priority *int `json:"priority,omitempty" validate:"omitempty,min=0"`
}

// WarehouseResponse represents a warehouse in API responses.
// @Description Warehouse information returned by the API
type WarehouseResponse struct {
	// ID is the unique warehouse identifier
	ID string `json:"id"`
	// Name is the warehouse display name
	Name string `json:"name"`
	// Code is the unique warehouse code
	Code string `json:"code"`
	// Address is the warehouse physical address
	Address WarehouseAddress `json:"address"`
	// IsActive indicates if the warehouse is operational
	IsActive bool `json:"is_active"`
	// Priority is used for stock allocation
	Priority int `json:"priority"`
	// TotalProducts is the count of distinct products in this warehouse
	TotalProducts int `json:"total_products"`
	// CreatedAt is when the warehouse was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the warehouse was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ListWarehousesRequest represents query parameters for listing warehouses.
type ListWarehousesRequest struct {
	PaginationRequest
	// ActiveOnly returns only active warehouses
	ActiveOnly bool `json:"active_only,omitempty"`
}

// ListWarehousesResponse represents the response for listing warehouses.
// @Description List of warehouses
type ListWarehousesResponse struct {
	// Warehouses is the list of warehouses
	Warehouses []WarehouseResponse `json:"warehouses"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}