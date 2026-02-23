// file: internal/api/types/product.go
package types

import (
	"github.com/google/uuid"
)

// VariantAttribute represents a product variant attribute (size, color, etc.)
type VariantAttribute struct {
	Name  string `json:"name" validate:"required,min=1,max=50"`
	Value string `json:"value" validate:"required,min=1,max=100"`
}

// Product represents a product in the inventory system
type Product struct {
	ID          uuid.UUID          `json:"id"`
	SKU         string             `json:"sku"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Category    string             `json:"category,omitempty"`
	Variants    []VariantAttribute `json:"variants,omitempty"`
	IsActive    bool               `json:"is_active"`
	LowStockThreshold int          `json:"low_stock_threshold"`
	AuditInfo
}

// CreateProductRequest represents the request body for creating a product
// @Description Request payload for creating a new product
type CreateProductRequest struct {
	// SKU is the unique stock keeping unit identifier
	SKU string `json:"sku" validate:"required,min=3,max=50,alphanumunicode"`
	// Name is the product display name
	Name string `json:"name" validate:"required,min=1,max=255"`
	// Description is an optional product description
	Description string `json:"description,omitempty" validate:"max=2000"`
	// Category is the product category
	Category string `json:"category,omitempty" validate:"max=100"`
	// LowStockThreshold triggers alerts when stock falls below this level
	LowStockThreshold int `json:"low_stock_threshold" validate:"min=0"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description       *string `json:"description,omitempty" validate:"omitempty,max=2000"`
	Category          *string `json:"category,omitempty" validate:"omitempty,max=100"`
	IsActive          *bool   `json:"is_active,omitempty"`
	LowStockThreshold *int    `json:"low_stock_threshold,omitempty" validate:"omitempty,min=0"`
}

// CreateVariantRequest represents the request body for creating a product variant
type CreateVariantRequest struct {
	// SKU is the unique identifier for this variant
	SKU string `json:"sku" validate:"required,min=3,max=50"`
	// Attributes define the variant properties (e.g., size: "L", color: "red")
	Attributes []VariantAttribute `json:"attributes" validate:"required,min=1,dive"`
}

// ProductResponse is the response wrapper for product operations
type ProductResponse struct {
	Product Product `json:"product"`
}

// ProductListResponse is the response for listing products
type ProductListResponse = PaginatedResponse[Product]

// ProductListParams contains query parameters for listing products
type ProductListParams struct {
	Pagination
	Category  string `json:"category,omitempty"`
	SKUPrefix string `json:"sku_prefix,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	Search    string `json:"search,omitempty"`
}