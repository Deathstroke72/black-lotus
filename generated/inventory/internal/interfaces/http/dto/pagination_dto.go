// file: internal/interfaces/http/dto/pagination_dto.go
package dto

// PaginationRequest represents pagination parameters for list endpoints.
type PaginationRequest struct {
	// Page number (1-indexed)
	Page int `json:"page" validate:"min=1"`
	// PageSize is the number of items per page
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse contains pagination metadata in list responses.
type PaginationResponse struct {
	// Page is the current page number
	Page int `json:"page"`
	// PageSize is the number of items per page
	PageSize int `json:"page_size"`
	// TotalItems is the total number of items across all pages
	TotalItems int64 `json:"total_items"`
	// TotalPages is the total number of pages
	TotalPages int `json:"total_pages"`
	// HasNext indicates if there are more pages
	HasNext bool `json:"has_next"`
	// HasPrev indicates if there are previous pages
	HasPrev bool `json:"has_prev"`
}

// DefaultPage is the default page number
const DefaultPage = 1

// DefaultPageSize is the default number of items per page
const DefaultPageSize = 20

// MaxPageSize is the maximum allowed page size
const MaxPageSize = 100