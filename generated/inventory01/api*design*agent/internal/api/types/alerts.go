// file: internal/api/types/alerts.go
package types

import (
	"github.com/google/uuid"
)

// AlertSeverity represents the severity of a low-stock alert
type AlertSeverity string

const (
	AlertSeverityWarning  AlertSeverity = "WARNING"  // Below threshold
	AlertSeverityCritical AlertSeverity = "CRITICAL" // Below 25% of threshold
	AlertSeverityOutOfStock AlertSeverity = "OUT_OF_STOCK"
)

// LowStockAlert represents a low-stock alert
type LowStockAlert struct {
	ID              uuid.UUID     `json:"id"`
	ProductID       uuid.UUID     `json:"product_id"`
	ProductSKU      string        `json:"product_sku"`
	ProductName     string        `json:"product_name"`
	WarehouseID     *uuid.UUID    `json:"warehouse_id,omitempty"` // Null for aggregate alerts
	WarehouseCode   *string       `json:"warehouse_code,omitempty"`
	CurrentQuantity int           `json:"current_quantity"`
	Threshold       int           `json:"threshold"`
	Severity        AlertSeverity `json:"severity"`
	CreatedAt       Timestamp     `json:"created_at"`
	AcknowledgedAt  *Timestamp    `json:"acknowledged_at,omitempty"`
	AcknowledgedBy  *string       `json:"acknowledged_by,omitempty"`
}

// LowStockAlertParams contains query parameters for alert queries
type LowStockAlertParams struct {
	Pagination
	Severity        *AlertSeverity `json:"severity,omitempty"`
	WarehouseID     *uuid.UUID     `json:"warehouse_id,omitempty"`
	Category        *string        `json:"category,omitempty"`
	Unacknowledged  bool           `json:"unacknowledged,omitempty"`
}

// LowStockAlertListResponse is the response for listing alerts
type LowStockAlertListResponse = PaginatedResponse[LowStockAlert]

// SetThresholdsRequest sets low-stock thresholds for a product
type SetThresholdsRequest struct {
	// GlobalThreshold applies to aggregated stock
	GlobalThreshold int `json:"global_threshold" validate:"min=0"`
	// WarehouseThresholds are warehouse-specific overrides
	WarehouseThresholds map[uuid.UUID]int `json:"warehouse_thresholds,omitempty"`
}

// ThresholdsResponse shows current thresholds for a product
type ThresholdsResponse struct {
	ProductID           uuid.UUID         `json:"product_id"`
	GlobalThreshold     int               `json:"global_threshold"`
	WarehouseThresholds map[uuid.UUID]int `json:"warehouse_thresholds"`
}