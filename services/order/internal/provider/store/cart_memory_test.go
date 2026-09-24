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

func TestMemoryCartStore_Find_NotFound(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryCartStore()
	_, err := s.Find(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Equal(t, domain.ErrCartNotFound, err)
}

func TestMemoryCartStore_UpdateAndFind(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryCartStore()
	customerID := uuid.New()
	cart := domain.Cart{RestaurantID: uuid.New(), SubtotalMinor: 100}

	err := s.Update(context.Background(), customerID, cart)
	assert.NoError(t, err)
	found, err := s.Find(context.Background(), customerID)
	assert.NoError(t, err)
	assert.Equal(t, cart, found)
}

func TestMemoryCartStore_UpdateClearsCart(t *testing.T) {
	t.Parallel()
	s := store.NewMemoryCartStore()
	customerID := uuid.New()

	s.Update(context.Background(), customerID, domain.Cart{RestaurantID: uuid.New(), SubtotalMinor: 100})

	newCart := domain.Cart{RestaurantID: uuid.New(), SubtotalMinor: 0}

	err := s.Update(context.Background(), customerID, newCart)
	assert.NoError(t, err)

	found, err := s.Find(context.Background(), customerID)
	assert.NoError(t, err)
	assert.Equal(t, newCart, found)
}

func TestMemoryCartStore_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	s := store.NewMemoryCartStore()
	customerID := uuid.New()

	for i := 0; i < 100; i++ {
		wg.Go(func() {
			err := s.Update(context.Background(), customerID, domain.Cart{RestaurantID: uuid.New(), SubtotalMinor: 100})
			assert.NoError(t, err)
		})
	}

	wg.Wait()

}
