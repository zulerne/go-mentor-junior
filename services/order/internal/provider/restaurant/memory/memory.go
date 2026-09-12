package memory

import (
	"context"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type RestaurantProvider struct {
	data map[string]map[string]domain.MenuItem
}

func NewRestaurantProvider() *RestaurantProvider {
	menuItems := map[string]map[string]domain.MenuItem{}
	menuItems["restaurant-1"] = map[string]domain.MenuItem{
		"menu-item-1": {
			ID:          "menu-item-1",
			Name:        "Lasagna",
			Description: "Delicious lasagna pasta",
			PriceMinor:  100,
			Currency:    "USD",
			Available:   true,
		},
		"menu-item-2": {
			ID:          "menu-item-2",
			Name:        "Spaghetti",
			Description: "Delicious spaghetti pasta",
			PriceMinor:  80,
			Currency:    "USD",
			Available:   true,
		},
	}
	return &RestaurantProvider{data: menuItems}
}

func (r *RestaurantProvider) Find(_ context.Context, restaurantID, menuItemID string) (domain.MenuItem, error) {
	items, ok := r.data[restaurantID]
	if !ok {
		return domain.MenuItem{}, nil
	}
	item, ok := items[menuItemID]
	if !ok {
		return domain.MenuItem{}, nil
	}
	return item, nil
}
