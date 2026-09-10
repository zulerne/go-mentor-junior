package memory

import (
	"context"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type RestaurantProvider struct {
	// restaurant_id -> menu_item_id
	data map[string]map[string]domain.MenuItem
}

func NewRestaurantProvider() *RestaurantProvider {
	menuItems := map[string]map[string]domain.MenuItem{}
	menuItems["restaurant-1"] = map[string]domain.MenuItem{
		"menu-item-1": {Id: "menu-item-1", Name: "Lasagna", Description: "Delicious lasagna pasta", PriceMinor: 100, Currency: "USD", Available: true},
		"menu-item-2": {Id: "menu-item-2", Name: "Spaghetti", Description: "Delicious spaghetti pasta", PriceMinor: 80, Currency: "USD", Available: true},
	}
	return &RestaurantProvider{data: menuItems}
}

func (r *RestaurantProvider) Find(ctx context.Context, restaurantId, menuItemId string) (domain.MenuItem, error) {
	items, ok := r.data[restaurantId]
	if !ok {
		return domain.MenuItem{}, nil
	}
	item, ok := items[menuItemId]
	if !ok {
		return domain.MenuItem{}, nil
	}
	return item, nil
}
