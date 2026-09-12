package memory

import (
	"context"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type CartStore struct {
	// data is a map of cart IDs to cart data
	data map[string]domain.Cart
}

// AddItem(ctx context.Context, restaurantId, menuItemId string, quantity int, instructions string) error

func NewCartStore() *CartStore {
	data := make(map[string]domain.Cart)

	restaurantID := "restaurant_id_1"
	currency := "USD"

	data["test"] = domain.Cart{
		RestaurantID: &restaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     "menu_item_id_1",
				Name:           "item_name",
				UnitPriceMinor: 100,
				Currency:       currency,
				Quantity:       1,
				Instructions:   "instructions",
			},
		},
		SubtotalMinor: 100,
		Currency:      &currency,
	}

	return &CartStore{
		data: data,
	}
}

func (s *CartStore) AddItem(
	_ context.Context,
	restaurantID, customerID string,
	menuItem domain.MenuItem,
	quantity int,
	instructions string,
) error {
	cart, ok := s.data[customerID]
	if !ok {
		cart = domain.Cart{
			RestaurantID:  &restaurantID,
			Items:         []domain.CartItem{},
			SubtotalMinor: 0,
			Currency:      &menuItem.Currency,
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

func (s *CartStore) FindCart(_ context.Context, customerID string) (domain.Cart, error) {
	cart, ok := s.data[customerID]
	if !ok {
		return domain.Cart{}, nil
	}
	return cart, nil
}
