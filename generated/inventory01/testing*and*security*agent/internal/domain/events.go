// file: internal/domain/events.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

type StockReservedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
}

type StockReleasedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	Reason        string    `json:"reason"`
	Timestamp     time.Time `json:"timestamp"`
}

type StockFulfilledEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	ProductID     uuid.UUID `json:"product_id"`
	WarehouseID   uuid.UUID `json:"warehouse_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	Timestamp     time.Time `json:"timestamp"`
}

type LowStockAlertEvent struct {
	ProductID         uuid.UUID `json:"product_id"`
	WarehouseID       uuid.UUID `json:"warehouse_id"`
	CurrentQuantity   int       `json:"current_quantity"`
	ReservedQuantity  int       `json:"reserved_quantity"`
	AvailableQuantity int       `json:"available_quantity"`
	ReorderPoint      int       `json:"reorder_point"`
	Timestamp         time.Time `json:"timestamp"`
}