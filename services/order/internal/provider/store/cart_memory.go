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

	restaurantID := uuid.MustParse("restaurant_id_1")
	currency := "USD"

	data[uuid.MustParse("test")] = domain.Cart{
		RestaurantID: restaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     uuid.MustParse("menu_item_id_1"),
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

func (s *MemoryCartStore) AddItem(
	_ context.Context,
	restaurantID uuid.UUID,
	customerID uuid.UUID,
	menuItem domain.MenuItem,
	quantity int,
	instructions string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.data[customerID]
	if !ok {
		cart = domain.Cart{
			RestaurantID:  restaurantID,
			Items:         []domain.CartItem{},
			SubtotalMinor: 0,
			Currency:      menuItem.Currency,
		}
	}
	s.data[customerID] = cart

	cart.Items = append(cart.Items, domain.CartItem{
		MenuItemID:   menuItem.ID,
		Quantity:     int32(quantity),
		Instructions: instructions,
	})
	cart.SubtotalMinor += int64(quantity) * menuItem.PriceMinor
	s.data[customerID] = cart
	return nil
}

func (s *MemoryCartStore) FindCart(_ context.Context, customerID uuid.UUID) (domain.Cart, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.data[customerID]
	if !ok {
		return domain.Cart{}, nil
	}

	cart.Items = slices.Clone(cart.Items)
	return cart, nil
}
