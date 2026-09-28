package store_test

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
	"github.com/zulerne/go-mentor-junior/order/internal/provider/store"
)

func TestMemoryOrderStore_Find_NotFound(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()
	_, err := s.Find(context.Background(), uuid.New())
	require.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestMemoryOrderStore_UpdateAndFind(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()
	order := domain.Order{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		Status:     domain.Pending,
	}
	updated, err := s.Update(context.Background(), order)
	require.NoError(t, err)
	assert.False(t, updated.UpdatedAt.IsZero(), "UpdatedAt should be set by store")
	found, err := s.Find(context.Background(), order.ID)
	require.NoError(t, err)
	assert.Equal(t, updated, found)
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
	updated, err := s.Update(context.Background(), order)
	require.NoError(t, err)
	found, err := s.FindByRestaurant(context.Background(), order.RestaurantID)
	require.NoError(t, err)
	assert.Equal(t, []domain.Order{updated}, found)
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
	updated, err := s.Update(context.Background(), order)
	require.NoError(t, err)
	found, err := s.GetAllOrders(context.Background(), order.CustomerID)
	require.NoError(t, err)
	assert.Equal(t, []domain.Order{updated}, found)
}

func TestMemoryOrderStore_GetAllOrders_Empty(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryOrderStore()

	found, err := s.GetAllOrders(context.Background(), uuid.New())
	require.NoError(t, err)
	var expected []domain.Order
	assert.Equal(t, expected, found)
}

func TestMemoryOrderStore_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	s := store.NewMemoryOrderStore()

	for range 100 {
		wg.Go(func() {
			order := domain.Order{
				ID:           uuid.New(),
				CustomerID:   uuid.New(),
				RestaurantID: uuid.New(),
				Status:       domain.Pending,
			}
			_, err := s.Update(context.Background(), order)
			require.NoError(t, err)
		})
	}

	wg.Wait()
}
