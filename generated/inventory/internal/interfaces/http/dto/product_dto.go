// file: internal/interfaces/http/dto/product_dto.go
package dto

import "time"

// ProductVariant represents a product variant (size, color combination).
type ProductVariant struct {
	// Size of the product variant (e.g., "S", "M", "L", "XL")
	Size string `json:"size,omitempty"`
	// Color of the product variant (e.g., "Red", "Blue")
	Color string `json:"color,omitempty"`
	// SKU is the unique stock keeping unit for this variant
	SKU string `json:"sku" validate:"required,min=1,max=100"`
}

// CreateProductRequest represents the request body for creating a product.
// @Description Request payload for creating a new product
type CreateProductRequest struct {
	// Name is the product display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Description is the product description
	Description string `json:"description,omitempty" validate:"max=2000"`
	// BaseSKU is the base SKU for the product (variants will extend this)
	BaseSKU string `json:"base_sku" validate:"required,min=1,max=100"`
	// Category is the product category
	Category string `json:"category,omitempty" validate:"max=100"`
	// Variants are the product variants (size, color combinations)
	Variants []ProductVariant `json:"variants,omitempty" validate:"dive"`
	// LowStockThreshold is the quantity below which low-stock alerts trigger
	LowStockThreshold int `json:"low_stock_threshold" validate:"min=0"`
	// Metadata contains additional product attributes
	Metadata map[string]string `json:"metadata,omitempty"`
}

// UpdateProductRequest represents the request body for updating a product.
// @Description Request payload for updating an existing product
type UpdateProductRequest struct {
	// Name is the product display name
	Name *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	// Description is the product description
	Description *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	// Category is the product category
	Category *string `json:"category,omitempty" validate:"omitempty,max=100"`
	// LowStockThreshold is the quantity below which low-stock alerts trigger
	LowStockThreshold *int `json:"low_stock_threshold,omitempty" validate:"omitempty,min=0"`
	// Metadata contains additional product attributes
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ProductResponse represents a product in API responses.
// @Description Product information returned by the API
type ProductResponse struct {
	// ID is the unique product identifier
	ID string `json:"id"`
	// Name is the product display name
	Name string `json:"name"`
	// Description is the product description
	Description string `json:"description,omitempty"`
	// BaseSKU is the base SKU for the product
	BaseSKU string `json:"base_sku"`
	// Category is the product category
	Category string `json:"category,omitempty"`
	// Variants are the product variants
	Variants []ProductVariant `json:"variants,omitempty"`
	// LowStockThreshold is the quantity below which low-stock alerts trigger
	LowStockThreshold int `json:"low_stock_threshold"`
	// TotalStock is the aggregated stock across all warehouses
	TotalStock int `json:"total_stock"`
	// TotalReserved is the aggregated reserved quantity
	TotalReserved int `json:"total_reserved"`
	// AvailableStock is TotalStock minus TotalReserved
	AvailableStock int `json:"available_stock"`
	// Metadata contains additional product attributes
	Metadata map[string]string `json:"metadata,omitempty"`
	// CreatedAt is when the product was created
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is when the product was last updated
	UpdatedAt time.Time `json:"updated_at"`
}

// ListProductsRequest represents query parameters for listing products.
type ListProductsRequest struct {
	PaginationRequest
	// Category filters by product category
	Category string `json:"category,omitempty"`
	// Search performs a text search on name and description
	Search string `json:"search,omitempty"`
	// LowStockOnly returns only products with low stock
	LowStockOnly bool `json:"low_stock_only,omitempty"`
}

// ListProductsResponse represents the response for listing products.
// @Description Paginated list of products
type ListProductsResponse struct {
	// Products is the list of products
	Products []ProductResponse `json:"products"`
	// Pagination contains pagination metadata
	Pagination PaginationResponse `json:"pagination"`
}