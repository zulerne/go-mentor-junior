package store_test

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store"
)

func TestMemoryOrderStore_Find_NotFound(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()
	_, err := s.Find(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, domain.ErrOrderNotFound, err)
}

func TestMemoryOrderStore_UpdateAndFind(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()
	order := domain.Order{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		Status:     domain.Pending,
	}
	err := s.Update(context.Background(), order)
	assert.NoError(t, err)
	found, err := s.Find(context.Background(), order.ID)
	assert.NoError(t, err)
	assert.Equal(t, order, found)
}

func TestMemoryOrderStore_FindByRestaurant(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()
	order := domain.Order{
		ID:           uuid.New(),
		CustomerID:   uuid.New(),
		RestaurantID: uuid.New(),
		Status:       domain.Pending,
	}
	err := s.Update(context.Background(), order)
	assert.NoError(t, err)
	found, err := s.FindByRestaurant(context.Background(), order.RestaurantID)
	assert.NoError(t, err)
	assert.Equal(t, []domain.Order{order}, found)
}

func TestMemoryOrderStore_GetAllOrders(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()
	order := domain.Order{
		ID:           uuid.New(),
		CustomerID:   uuid.New(),
		RestaurantID: uuid.New(),
		Status:       domain.Pending,
	}
	err := s.Update(context.Background(), order)
	assert.NoError(t, err)
	found, err := s.GetAllOrders(context.Background(), order.CustomerID)
	assert.NoError(t, err)
	assert.Equal(t, []domain.Order{order}, found)
}

func TestMemoryOrderStore_GetAllOrders_Empty(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()

	found, err := s.GetAllOrders(context.Background(), uuid.New())
	assert.NoError(t, err)
	var expected []domain.Order
	assert.Equal(t, expected, found)
}

func TestMemoryOrderStore_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	s := store.NewMemoryOrderStore()

	for i := 0; i < 100; i++ {
		wg.Go(func() {
			order := domain.Order{
				ID:           uuid.New(),
				CustomerID:   uuid.New(),
				RestaurantID: uuid.New(),
				Status:       domain.Pending,
			}
			err := s.Update(context.Background(), order)
			assert.NoError(t, err)
		})
	}

	wg.Wait()
}
