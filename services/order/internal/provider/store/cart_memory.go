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

	for i := range cart.Items {
		if cart.Items[i].MenuItemID == menuItem.ID {
			oldPrice, oldQuantity := cart.Items[i].UnitPriceMinor, cart.Items[i].Quantity
			cart.Items[i].UnitPriceMinor = menuItem.PriceMinor
			cart.Items[i].Quantity = int32(quantity)
			cart.SubtotalMinor -= int64(oldQuantity) * oldPrice
			cart.SubtotalMinor += int64(quantity) * menuItem.PriceMinor

			cart.Items[i].Name = menuItem.Name
			cart.Items[i].Instructions = instructions
			s.data[customerID] = cart
			return nil
		}
	}

	cart.Items = append(cart.Items, domain.CartItem{
		MenuItemID:     menuItem.ID,
		Name:           menuItem.Name,
		Quantity:       int32(quantity),
		UnitPriceMinor: menuItem.PriceMinor,
		Currency:       menuItem.Currency,
		Instructions:   instructions,
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
