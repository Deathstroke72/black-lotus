// file: internal/domain/entity/product.go
package entity

import (
	"errors"
	"time"
)

// ProductVariant represents a specific variant of a product (e.g., size, color)
type ProductVariant struct {
	Size  string
	Color string
}

// Product represents a product in the inventory system
type Product struct {
	ID          string
	SKU         string
	Name        string
	Description string
	Variant     ProductVariant
	Category    string
	MinStock    int // Threshold for low-stock alerts
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Validation errors
var (
	ErrProductIDRequired   = errors.New("product ID is required")
	ErrProductSKURequired  = errors.New("product SKU is required")
	ErrProductNameRequired = errors.New("product name is required")
	ErrMinStockNegative    = errors.New("minimum stock cannot be negative")
	ErrProductDeleted      = errors.New("product has been deleted")
)

// NewProduct creates a new Product with validation
func NewProduct(id, sku, name, description, category string, variant ProductVariant, minStock int) (*Product, error) {
	if id == "" {
		return nil, ErrProductIDRequired
	}
	if sku == "" {
		return nil, ErrProductSKURequired
	}
	if name == "" {
		return nil, ErrProductNameRequired
	}
	if minStock < 0 {
		return nil, ErrMinStockNegative
	}

	now := time.Now().UTC()
	return &Product{
		ID:          id,
		SKU:         sku,
		Name:        name,
		Description: description,
		Variant:     variant,
		Category:    category,
		MinStock:    minStock,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update modifies product details
func (p *Product) Update(name, description, category string, variant ProductVariant, minStock int) error {
	if p.DeletedAt != nil {
		return ErrProductDeleted
	}
	if name == "" {
		return ErrProductNameRequired
	}
	if minStock < 0 {
		return ErrMinStockNegative
	}

	p.Name = name
	p.Description = description
	p.Category = category
	p.Variant = variant
	p.MinStock = minStock
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// SoftDelete marks the product as deleted
func (p *Product) SoftDelete() error {
	if p.DeletedAt != nil {
		return ErrProductDeleted
	}
	now := time.Now().UTC()
	p.DeletedAt = &now
	p.IsActive = false
	p.UpdatedAt = now
	return nil
}

// IsDeleted returns true if the product has been soft deleted
func (p *Product) IsDeleted() bool {
	return p.DeletedAt != nil
}