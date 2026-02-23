// file: internal/domain/repository/product_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// ProductFilter defines filtering options for product queries
type ProductFilter struct {
	SKU      *string
	Name     *string
	Category *string
	IsActive *bool
	Limit    int
	Offset   int
}

// ProductRepository defines the interface for product persistence
type ProductRepository interface {
	// Create persists a new product
	Create(ctx context.Context, product *entity.Product) error

	// GetByID retrieves a product by its ID
	GetByID(ctx context.Context, id string) (*entity.Product, error)

	// GetBySKU retrieves a product by its SKU
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)

	// List retrieves products with optional filtering
	List(ctx context.Context, filter ProductFilter) ([]*entity.Product, int, error)

	// Update persists changes to an existing product
	Update(ctx context.Context, product *entity.Product) error

	// Delete soft deletes a product
	Delete(ctx context.Context, id string) error

	// ExistsBySKU checks if a product with the given SKU exists
	ExistsBySKU(ctx context.Context, sku string) (bool, error)
}