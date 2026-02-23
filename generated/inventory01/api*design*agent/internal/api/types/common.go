// file: internal/api/types/common.go
package types

import (
	"time"
)

// Pagination represents pagination parameters for list endpoints
type Pagination struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginatedResponse wraps list responses with pagination metadata
type PaginatedResponse[T any] struct {
	Data       []T              `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

// PaginationMeta contains pagination metadata in responses
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// Timestamp is a wrapper for consistent time formatting
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).UTC().Format(time.RFC3339) + `"`), nil
}

// AuditInfo contains common audit fields
type AuditInfo struct {
	CreatedAt Timestamp  `json:"created_at"`
	UpdatedAt Timestamp  `json:"updated_at"`
	CreatedBy string     `json:"created_by,omitempty"`
	UpdatedBy string     `json:"updated_by,omitempty"`
}