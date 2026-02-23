// file: internal/domain/repository/stock_movement_repository.go
package repository

import (
	"context"
	"time"

	"github.com/inventory-service/internal/domain/entity"
)

// StockMovementFilter defines filtering options for stock movement queries
type StockMovementFilter struct {
	StockItemID   *string
	MovementType  *entity.MovementType
	ReferenceID   *string
	ReferenceType *string
	StartDate     *time.Time
	EndDate       *time.Time
	Limit         int
	Offset        int
}

// StockMovementRepository defines the interface for stock movement persistence
type StockMovementRepository interface {
	// Create persists a new stock movement record
	Create(ctx context.Context, movement *entity.StockMovement) error

	// GetByID retrieves a stock movement by its ID
	GetByID(ctx context.Context, id string) (*entity.StockMovement, error)

	// List retrieves stock movements with optional filtering
	List(ctx context.Context, filter StockMovementFilter) ([]*entity.StockMovement, int, error)

	// GetByStockItem retrieves all movements for a specific stock item
	GetByStockItem(ctx context.Context, stockItemID string, limit, offset int) ([]*entity.StockMovement, int, error)

	// GetByReference retrieves movements by reference (e.g., order ID)
	GetByReference(ctx context.Context, referenceID, referenceType string) ([]*entity.StockMovement, error)
}