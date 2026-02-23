// file: internal/testutil/mocks.go
package testutil

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"inventory/internal/domain"
	"inventory/internal/repository"
)

// MockProductRepository is a mock implementation of ProductRepository
type MockProductRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{}
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) List(ctx context.Context, filter repository.ProductFilter) ([]*domain.Product, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Product), args.Int(1), args.Error(2)
}

// MockStockRepository is a mock implementation of StockRepository
type MockStockRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockStockRepository() *MockStockRepository {
	return &MockStockRepository{}
}

func (m *MockStockRepository) GetByProductAndWarehouse(ctx context.Context, productID, warehouseID uuid.UUID) (*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID, warehouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StockItem), args.Error(1)
}

func (m *MockStockRepository) GetByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockItem), args.Error(1)
}

func (m *MockStockRepository) GetByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, warehouseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockItem), args.Error(1)
}

func (m *MockStockRepository) Create(ctx context.Context, item *domain.StockItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockStockRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity, reserved int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity, reserved)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicReserve(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicRelease(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicFulfill(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) AtomicReplenish(ctx context.Context, id uuid.UUID, quantity int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockStockRepository) GetLowStock(ctx context.Context, threshold int) ([]*domain.StockItem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockItem), args.Error(1)
}

// MockReservationRepository is a mock implementation of ReservationRepository
type MockReservationRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockReservationRepository() *MockReservationRepository {
	return &MockReservationRepository{}
}

func (m *MockReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, reservation)
	return args.Error(0)
}

func (m *MockReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

func (m *MockReservationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReservationStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockReservationRepository) List(ctx context.Context, filter repository.ReservationFilter) ([]*domain.Reservation, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Reservation), args.Int(1), args.Error(2)
}

func (m *MockReservationRepository) GetExpired(ctx context.Context, before time.Time) ([]*domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, before)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Reservation), args.Error(1)
}

// MockStockMovementRepository is a mock implementation of StockMovementRepository
type MockStockMovementRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockStockMovementRepository() *MockStockMovementRepository {
	return &MockStockMovementRepository{}
}

func (m *MockStockMovementRepository) Create(ctx context.Context, movement *domain.StockMovement) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, movement)
	return args.Error(0)
}

func (m *MockStockMovementRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.StockMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StockMovement), args.Error(1)
}

func (m *MockStockMovementRepository) List(ctx context.Context, filter repository.MovementFilter) ([]*domain.StockMovement, int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.StockMovement), args.Int(1), args.Error(2)
}

func (m *MockStockMovementRepository) GetByProduct(ctx context.Context, productID uuid.UUID, limit int) ([]*domain.StockMovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.StockMovement), args.Error(1)
}

// MockWarehouseRepository is a mock implementation of WarehouseRepository
type MockWarehouseRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockWarehouseRepository() *MockWarehouseRepository {
	return &MockWarehouseRepository{}
}

func (m *MockWarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) GetByCode(ctx context.Context, code string) (*domain.Warehouse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Warehouse), args.Error(1)
}

func (m *MockWarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, warehouse)
	return args.Error(0)
}

func (m *MockWarehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockWarehouseRepository) List(ctx context.Context, activeOnly bool) ([]*domain.Warehouse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, activeOnly)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Warehouse), args.Error(1)
}

// MockCacheRepository is a mock implementation of CacheRepository
type MockCacheRepository struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockCacheRepository() *MockCacheRepository {
	return &MockCacheRepository{}
}

func (m *MockCacheRepository) GetStockLevel(ctx context.Context, productID uuid.UUID) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	args := m.Called(ctx, productID)
	return args.Int(0), args.Error(1)
}

func (m *MockCacheRepository) SetStockLevel(ctx context.Context, productID uuid.UUID, quantity int, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, productID, quantity, ttl)
	return args.Error(0)
}

func (m *MockCacheRepository) InvalidateStockLevel(ctx context.Context, productID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, productID)
	return args.Error(0)
}

// MockEventPublisher is a mock implementation of EventPublisher
type MockEventPublisher struct {
	mock.Mock
	mu sync.RWMutex
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

func (m *MockEventPublisher) PublishStockReserved(ctx context.Context, event domain.StockReservedEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishStockReleased(ctx context.Context, event domain.StockReleasedEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishLowStockAlert(ctx context.Context, event domain.LowStockAlertEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventPublisher) PublishStockFulfilled(ctx context.Context, event domain.StockFulfilledEvent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockTransactionManager is a mock implementation of TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

func NewMockTransactionManager() *MockTransactionManager {
	return &MockTransactionManager{}
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	// Execute the function to test transaction behavior
	if args.Bool(1) {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return args.Error(0)
}