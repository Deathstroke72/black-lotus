// file: internal/events/consumed_events.go
package events

import (
	"encoding/json"
	"time"
)

// --- Consumed Events from Order Service ---

// OrderCreatedEvent is consumed when a new order is created
type OrderCreatedEvent struct {
	BaseEvent
	Payload OrderCreatedPayload `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID       string           `json:"order_id"`
	CustomerID    string           `json:"customer_id"`
	Items         []OrderItem      `json:"items"`
	ShippingAddr  Address          `json:"shipping_address"`
	OrderedAt     time.Time        `json:"ordered_at"`
	Priority      string           `json:"priority,omitempty"`
}

type OrderItem struct {
	ProductID   string  `json:"product_id"`
	VariantID   string  `json:"variant_id,omitempty"`
	SKU         string  `json:"sku"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// OrderCancelledEvent is consumed when an order is cancelled
type OrderCancelledEvent struct {
	BaseEvent
	Payload OrderCancelledPayload `json:"payload"`
}

type OrderCancelledPayload struct {
	OrderID      string    `json:"order_id"`
	CustomerID   string    `json:"customer_id"`
	Reason       string    `json:"reason"`
	CancelledBy  string    `json:"cancelled_by"`
	CancelledAt  time.Time `json:"cancelled_at"`
}

// OrderFulfilledEvent is consumed when an order is fulfilled/shipped
type OrderFulfilledEvent struct {
	BaseEvent
	Payload OrderFulfilledPayload `json:"payload"`
}

type OrderFulfilledPayload struct {
	OrderID       string             `json:"order_id"`
	CustomerID    string             `json:"customer_id"`
	Items         []FulfilledItem    `json:"items"`
	ShipmentID    string             `json:"shipment_id"`
	FulfilledAt   time.Time          `json:"fulfilled_at"`
	WarehouseID   string             `json:"warehouse_id"`
}

type FulfilledItem struct {
	ProductID   string `json:"product_id"`
	VariantID   string `json:"variant_id,omitempty"`
	SKU         string `json:"sku"`
	Quantity    int    `json:"quantity"`
}

// OrderModifiedEvent is consumed when an order is modified
type OrderModifiedEvent struct {
	BaseEvent
	Payload OrderModifiedPayload `json:"payload"`
}

type OrderModifiedPayload struct {
	OrderID       string          `json:"order_id"`
	CustomerID    string          `json:"customer_id"`
	OldItems      []OrderItem     `json:"old_items"`
	NewItems      []OrderItem     `json:"new_items"`
	ModifiedBy    string          `json:"modified_by"`
	ModifiedAt    time.Time       `json:"modified_at"`
	Reason        string          `json:"reason,omitempty"`
}

// Deserialize methods for consumed events
func DeserializeOrderCreatedEvent(data []byte) (*OrderCreatedEvent, error) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DeserializeOrderCancelledEvent(data []byte) (*OrderCancelledEvent, error) {
	var event OrderCancelledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DeserializeOrderFulfilledEvent(data []byte) (*OrderFulfilledEvent, error) {
	var event OrderFulfilledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

func DeserializeOrderModifiedEvent(data []byte) (*OrderModifiedEvent, error) {
	var event OrderModifiedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}