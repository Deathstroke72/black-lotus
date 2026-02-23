// file: internal/domain/entity/low_stock_alert.go
package entity

import (
	"errors"
	"time"
)

// AlertStatus represents the current state of a low stock alert
type AlertStatus string

const (
	AlertStatusActive      AlertStatus = "ACTIVE"
	AlertStatusAcknowledged AlertStatus = "ACKNOWLEDGED"
	AlertStatusResolved    AlertStatus = "RESOLVED"
)

// LowStockAlert represents an alert for low stock levels
type LowStockAlert struct {
	ID              string
	StockItemID     string
	ProductID       string
	WarehouseID     string
	CurrentQuantity int
	ReorderPoint    int
	Status          AlertStatus
	AcknowledgedBy  *string
	AcknowledgedAt  *time.Time
	ResolvedAt      *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// LowStockAlert validation errors
var (
	ErrAlertIDRequired        = errors.New("alert ID is required")
	ErrAlertStockItemRequired = errors.New("stock item ID is required")
	ErrAlertProductRequired   = errors.New("product ID is required")
	ErrAlertWarehouseRequired = errors.New("warehouse ID is required")
	ErrAlertAlreadyResolved   = errors.New("alert has already been resolved")
)

// NewLowStockAlert creates a new LowStockAlert
func NewLowStockAlert(id, stockItemID, productID, warehouseID string, currentQuantity, reorderPoint int) (*LowStockAlert, error) {
	if id == "" {
		return nil, ErrAlertIDRequired
	}
	if stockItemID == "" {
		return nil, ErrAlertStockItemRequired
	}
	if productID == "" {
		return nil, ErrAlertProductRequired
	}
	if warehouseID == "" {
		return nil, ErrAlertWarehouseRequired
	}

	now := time.Now().UTC()
	return &LowStockAlert{
		ID:              id,
		StockItemID:     stockItemID,
		ProductID:       productID,
		WarehouseID:     warehouseID,
		CurrentQuantity: currentQuantity,
		ReorderPoint:    reorderPoint,
		Status:          AlertStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Acknowledge marks the alert as acknowledged
func (a *LowStockAlert) Acknowledge(userID string) error {
	if a.Status == AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	now := time.Now().UTC()
	a.Status = AlertStatusAcknowledged
	a.AcknowledgedBy = &userID
	a.AcknowledgedAt = &now
	a.UpdatedAt = now
	return nil
}

// Resolve marks the alert as resolved
func (a *LowStockAlert) Resolve() error {
	if a.Status == AlertStatusResolved {
		return ErrAlertAlreadyResolved
	}

	now := time.Now().UTC()
	a.Status = AlertStatusResolved
	a.ResolvedAt = &now
	a.UpdatedAt = now
	return nil
}