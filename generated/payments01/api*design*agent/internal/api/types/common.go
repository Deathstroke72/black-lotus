// file: internal/api/types/common.go
package types

import "time"

// Money represents a monetary amount with currency
// Currency follows ISO 4217 (e.g., USD, EUR, GBP)
// Amount is in the smallest currency unit (cents for USD)
type Money struct {
	Amount   int64  `json:"amount" validate:"required,gt=0"`
	Currency string `json:"currency" validate:"required,len=3,uppercase"`
}

// PaginationParams for list endpoints
type PaginationParams struct {
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" validate:"omitempty,min=0"`
	Cursor string `json:"cursor,omitempty"`
}

// PaginationMeta returned in list responses
type PaginationMeta struct {
	Total      int64  `json:"total"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

// Timestamp wrapper for consistent JSON formatting
type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}