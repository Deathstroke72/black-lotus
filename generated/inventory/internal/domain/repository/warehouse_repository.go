// file: internal/domain/repository/warehouse_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// WarehouseFilter defines filtering options for warehouse queries
type WarehouseFilter struct {
	Code     *string
	Name     *string
	IsActive *bool
	Limit    int
	Offset   int
}

// WarehouseRepository defines the interface for warehouse persistence
type WarehouseRepository interface {
	// Create persists a new warehouse
	Create(ctx context.Context, warehouse *entity.Warehouse) error

	// GetByID retrieves a warehouse by its ID
	GetByID(ctx context.Context, id string) (*entity.Warehouse, error)

	// GetByCode retrieves a warehouse by its code
	GetByCode(ctx context.Context, code string) (*entity.Warehouse, error)

	// List retrieves warehouses with optional filtering
	List(ctx context.Context, filter WarehouseFilter) ([]*entity.Warehouse, int, error)

	// Update persists changes to an existing warehouse
	Update(ctx context.Context, warehouse *entity.Warehouse) error

	// Delete soft deletes a warehouse
	Delete(ctx context.Context, id string) error

	// ExistsByCode checks if a warehouse with the given code exists
	ExistsByCode(ctx context.Context, code string) (bool, error)
}