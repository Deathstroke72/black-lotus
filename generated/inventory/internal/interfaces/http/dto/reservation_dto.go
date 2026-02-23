// file: internal/interfaces/http/dto/reservation_dto.go
package dto

import "time"

// ReservationItem represents a single item in a reservation request.
type ReservationItem struct {
	// ProductID is the product to reserve
	ProductID string `json:"product_id" validate:"required,uuid"`
	// VariantSKU is the specific variant SKU (optional)
	VariantSKU string `json:"variant_sku,omitempty" validate:"max=100"`
	// Quantity is the amount to reserve
	Quantity int `json:"quantity" validate:"required,min=1"`
	// PreferredWarehouseID is the preferred warehouse (optional)
	PreferredWarehouseID string `json:"preferred_warehouse_id,omitempty" validate:"omitempty,uuid"`
}

// CreateReservationRequest represents the request body for creating a reservation.
// @Description Request payload for reserving stock for an order
type CreateReservationRequest struct {
	// OrderID is the external order identifier
	OrderID string `json:"order_id" validate:"required,min=1,max=100"`
	// Items are the products and quantities to reserve
	Items []ReservationItem `json:"items" validate:"required,min=1,dive"`
	// ExpiresAt is when the reservation should expire (optional)
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Metadata contains additional reservation context
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ReservationItemResponse represents a reserved item in the response.
type ReservationItemResponse struct {
	// ProductID is the reserved product
	ProductID string `json:"product_id"`
	// ProductName is the product name
	ProductName string `json:"product_name"`
	// VariantSKU is the variant SKU
	VariantSKU string `json:"variant_sku,omitempty"`
	// Quantity is the reserved amount
	Quantity int `json:"quantity"`
	// WarehouseID is where the stock is reserved
	WarehouseID string `json:"warehouse_id"`
	// WarehouseName is the warehouse name
	WarehouseName string `json:"warehouse_name"`
	// StockItemID is the specific stock item
	StockItemID string `json:"stock_item_id"`
}

// ReservationResponse represents a reservation in API responses.
// @Description Reservation information returned by the API
type ReservationResponse struct {
	// ID is the unique reservation identifier
	ID string `json:"id"`
	// OrderID is the external order identifier
	OrderID string `json:"order_id"`
	// Status is the reservation status (pending, confirmed, released, fulfilled, expired)
	Status string `json:"status"`
	// Items are the reserved items
	Items []ReservationItemResponse `json:"items"`
	// ExpiresAt is when the reservation expires
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Metadata contains additional reservation context
	Metadata map[string]string `json:"metadata,omitempty"`
	// CreatedAt is when the reservation was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the reservation was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ReleaseReservationRequest represents the request body for releasing a reservation.
// @Description Request payload for releasing reserved stock
type ReleaseReservationRequest struct {
	// Reason is the reason for releasing (e.g., "order_cancelled", "timeout")
	Reason string `json:"reason" validate:"required,min=1,max=255"`
	// PartialItems allows releasing only specific items (optional)
	PartialItems []PartialReleaseItem `json:"partial_items,omitempty" validate:"omitempty,dive"`
}

// PartialReleaseItem specifies partial release of a reservation item.
type PartialReleaseItem struct {
	// ProductID is the product to partially release
	ProductID string `json:"product_id" validate:"required,uuid"`
	// VariantSKU is the variant SKU
	VariantSKU string `json:"variant_sku,omitempty"`
	// Quantity is the amount to release
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// FulfillReservationRequest represents the request body for fulfilling a reservation.
// @Description Request payload for fulfilling a reservation (decrementing stock)
type FulfillReservationRequest struct {
	// ShipmentID is the external shipment identifier (optional)
	ShipmentID string `json:"shipment_id,omitempty" validate:"max=100"`
	// FulfilledBy is the user or system that fulfilled the reservation
	FulfilledBy string `json:"fulfilled_by" validate:"required,max=255"`
	// Notes contains any fulfillment notes
	Notes string `json:"notes,omitempty" validate:"max=1000"`
}

// ListReservationsRequest represents query parameters for listing reservations.
type ListReservationsRequest struct {
	PaginationRequest
	// OrderID filters by order
	OrderID string `json:"order_id,omitempty"`
	// Status filters by reservation status
	Status string `json:"status,omitempty"`
	// ProductID filters by product
	ProductID string `json:"product_id,omitempty" validate:"omitempty,uuid"`
}

// ListReservationsResponse represents the response for listing reservations.
// @Description Paginated list of reservations
type ListReservationsResponse struct {
	// Reservations is the list of reservations
	Reservations []ReservationResponse `json:"reservations"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}

// Reservation status constants
const (
	ReservationStatusPending   = "pending"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusReleased  = "released"
	ReservationStatusFulfilled = "fulfilled"
	ReservationStatusExpired   = "expired"
)