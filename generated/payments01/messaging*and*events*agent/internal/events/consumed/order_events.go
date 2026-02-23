// file: internal/events/consumed/order_events.go
package consumed

import (
	"encoding/json"
	"time"
)

// Order Service Events that Payments service consumes

const (
	EventTypeOrderCreated   = "order.created"
	EventTypeOrderConfirmed = "order.confirmed"
	EventTypeOrderCanceled  = "order.canceled"
	EventTypeOrderCompleted = "order.completed"
)

// OrderCreatedEvent is consumed to initiate payment processing
type OrderCreatedEvent struct {
	EventID       string              `json:"event_id"`
	EventType     string              `json:"event_type"`
	EventVersion  string              `json:"event_version"`
	AggregateID   string              `json:"aggregate_id"`
	CorrelationID string              `json:"correlation_id"`
	Timestamp     time.Time           `json:"timestamp"`
	Payload       OrderCreatedPayload `json:"payload"`
}

type OrderCreatedPayload struct {
	OrderID          string           `json:"order_id"`
	CustomerID       string           `json:"customer_id"`
	TotalAmount      int64            `json:"total_amount"`
	Currency         string           `json:"currency"`
	Items            []OrderItem      `json:"items"`
	ShippingAddress  Address          `json:"shipping_address"`
	BillingAddress   Address          `json:"billing_address"`
	PaymentMethodID  string           `json:"payment_method_id,omitempty"`
	IdempotencyKey   string           `json:"idempotency_key"`
	RequireCapture   bool             `json:"require_capture"` // true = authorize only, false = capture immediately
	Metadata         map[string]string `json:"metadata,omitempty"`
}

type OrderItem struct {
	ItemID      string `json:"item_id"`
	ProductID   string `json:"product_id"`
	Name        string `json:"name"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int64  `json:"unit_price"`
	TotalPrice  int64  `json:"total_price"`
}

type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// OrderCanceledEvent is consumed to cancel/refund payment
type OrderCanceledEvent struct {
	EventID       string               `json:"event_id"`
	EventType     string               `json:"event_type"`
	EventVersion  string               `json:"event_version"`
	AggregateID   string               `json:"aggregate_id"`
	CorrelationID string               `json:"correlation_id"`
	Timestamp     time.Time            `json:"timestamp"`
	Payload       OrderCanceledPayload `json:"payload"`
}

type OrderCanceledPayload struct {
	OrderID       string    `json:"order_id"`
	CustomerID    string    `json:"customer_id"`
	Reason        string    `json:"reason"`
	CanceledBy    string    `json:"canceled_by"`
	RefundAmount  int64     `json:"refund_amount"` // 0 if no refund needed
	CanceledAt    time.Time `json:"canceled_at"`
}

// OrderConfirmedEvent is consumed to capture authorized payment
type OrderConfirmedEvent struct {
	EventID       string                `json:"event_id"`
	EventType     string                `json:"event_type"`
	EventVersion  string                `json:"event_version"`
	AggregateID   string                `json:"aggregate_id"`
	CorrelationID string                `json:"correlation_id"`
	Timestamp     time.Time             `json:"timestamp"`
	Payload       OrderConfirmedPayload `json:"payload"`
}

type OrderConfirmedPayload struct {
	OrderID      string    `json:"order_id"`
	CustomerID   string    `json:"customer_id"`
	ConfirmedAt  time.Time `json:"confirmed_at"`
}

// ParseOrderCreatedEvent parses raw JSON into OrderCreatedEvent
func ParseOrderCreatedEvent(data []byte) (*OrderCreatedEvent, error) {
	var event OrderCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// ParseOrderCanceledEvent parses raw JSON into OrderCanceledEvent
func ParseOrderCanceledEvent(data []byte) (*OrderCanceledEvent, error) {
	var event OrderCanceledEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// ParseOrderConfirmedEvent parses raw JSON into OrderConfirmedEvent
func ParseOrderConfirmedEvent(data []byte) (*OrderConfirmedEvent, error) {
	var event OrderConfirmedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}