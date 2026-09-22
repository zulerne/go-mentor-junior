package store

import (
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type MemoryCartStore struct {
	// data is a map of cart IDs to cart data
	data map[uuid.UUID]domain.Cart
	mu   sync.RWMutex
}

func NewMemoryCartStore() *MemoryCartStore {
	data := make(map[uuid.UUID]domain.Cart)

	// restId: fe70cb38-d10d-452c-8860-1af5936f7037
	// id1: 58185771-1f9f-4d71-9609-bdacea1deb2e
	// id2: 4c5344f5-33d7-4c05-bbbc-a0edc4178e7e
	// userId: ef55f77a-7738-426f-93a4-f78f2baf7970
	// orderId: 2f5efdc6-c936-4fdf-a681-1e21e73e6e71
	restaurantID := uuid.MustParse("fe70cb38-d10d-452c-8860-1af5936f7037")
	currency := "USD"
	userId := uuid.MustParse("ef55f77a-7738-426f-93a4-f78f2baf7970")
	data[userId] = domain.Cart{
		RestaurantID: restaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     uuid.MustParse("58185771-1f9f-4d71-9609-bdacea1deb2e"),
				Name:           "item_name",
				UnitPriceMinor: 100,
				Currency:       currency,
				Quantity:       1,
				Instructions:   "instructions",
			},
		},
		SubtotalMinor: 100,
		Currency:      currency,
	}

	return &MemoryCartStore{
		data: data,
	}
}

func (s *MemoryCartStore) Find(_ context.Context, customerID uuid.UUID) (domain.Cart, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.data[customerID]
	if !ok {
		return domain.Cart{}, domain.ErrCartNotFound
	}

	cart.Items = slices.Clone(cart.Items)
	return cart, nil
}

func (s *MemoryCartStore) Update(_ context.Context, customerID uuid.UUID, cart domain.Cart) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[customerID] = cart
	return nil
}
