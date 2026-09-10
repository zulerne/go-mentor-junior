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

	restaurantId := "restaurant_id_1"
	currency := "USD"

	data["test"] = domain.Cart{
		RestaurantID: &restaurantId,
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

func (s *CartStore) AddItem(ctx context.Context, restaurantId, customerId string, menuItem domain.MenuItem, quantity int, instructions string) error {
	cart, ok := s.data[customerId]
	if !ok {
		cart = domain.Cart{
			RestaurantID:  &restaurantId,
			Items:         []domain.CartItem{},
			SubtotalMinor: 0,
			Currency:      &menuItem.Currency,
		}
	}
	s.data[customerId] = cart

	cart.Items = append(cart.Items, domain.CartItem{
		MenuItemID:   menuItem.Id,
		Quantity:     int32(quantity),
		Instructions: instructions,
	})
	cart.SubtotalMinor += int64(quantity) * menuItem.PriceMinor
	s.data[customerId] = cart
	return nil
}

func (s *CartStore) FindCart(ctx context.Context, customerId string) (domain.Cart, error) {
	cart, ok := s.data[customerId]
	if !ok {
		return domain.Cart{}, nil
	}
	return cart, nil
}
