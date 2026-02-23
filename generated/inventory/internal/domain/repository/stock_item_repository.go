// file: internal/domain/repository/stock_item_repository.go
package repository

import (
	"context"

	"github.com/inventory-service/internal/domain/entity"
)

// StockItemFilter defines filtering options for stock item queries
type StockItemFilter struct {
	ProductID   *string
	WarehouseID *string
	LowStock    *bool // Filter items at or below reorder point
	Limit       int
	Offset      int
}

// AggregatedStock represents total stock for a product across warehouses
type AggregatedStock struct {
	ProductID        string
	TotalOnHand      int
	TotalReserved    int
	TotalAvailable   int
	WarehouseCount   int
	WarehouseDetails []WarehouseStockDetail
}

// WarehouseStockDetail represents stock in a specific warehouse
type WarehouseStockDetail struct {
	WarehouseID      string
	WarehouseName    string
	QuantityOnHand   int
	QuantityReserved int
	Available        int
}

// StockItemRepository defines the interface for stock item persistence
type StockItemRepository interface {
	// Create persists a new stock item
	Create(ctx context.Context, stockItem *entity.StockItem) error

	// GetByID retrieves a stock item by its ID
	GetByID(ctx context.Context, id string) (*entity.StockItem, error)

	// GetByProductAndWarehouse retrieves a stock item by product and warehouse
	GetByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (*entity.StockItem, error)

	// List retrieves stock items with optional filtering
	List(ctx context.Context, filter StockItemFilter) ([]*entity.StockItem, int, error)

	// Update persists changes to an existing stock item
	Update(ctx context.Context, stockItem *entity.StockItem) error

	// UpdateWithLock updates a stock item with optimistic locking
	UpdateWithLock(ctx context.Context, stockItem *entity.StockItem, expectedVersion int) error

	// GetAggregatedStock retrieves total stock for a product across all warehouses
	GetAggregatedStock(ctx context.Context, productID string) (*AggregatedStock, error)

	// GetLowStockItems retrieves all stock items at or below reorder point
	GetLowStockItems(ctx context.Context) ([]*entity.StockItem, error)

	// ExistsByProductAndWarehouse checks if a stock item exists for the given product and warehouse
	ExistsByProductAndWarehouse(ctx context.Context, productID, warehouseID string) (bool, error)
}