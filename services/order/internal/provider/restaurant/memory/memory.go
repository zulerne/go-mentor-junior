package memory

import (
	"context"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type Restaurant struct {
	MenuItems []domain.MenuItem
}

func NewRestaurant() *Restaurant {
	menuItems := []domain.MenuItem{}
	menuItems = append(menuItems, domain.MenuItem{Id: "menu-item-1", Name: "Lasagna", Description: "Delicious lasagna pasta", PriceMinor: 100, Currency: "USD", Available: true})
	menuItems = append(menuItems, domain.MenuItem{Id: "menu-item-2", Name: "Spaghetti", Description: "Delicious spaghetti pasta", PriceMinor: 80, Currency: "USD", Available: true})
	return &Restaurant{MenuItems: menuItems}
}

func (r *Restaurant) Find(ctx context.Context, menuItemId string) (domain.MenuItem, error) {
	for _, item := range r.MenuItems {
		if item.Id == menuItemId {
			return item, nil
		}
	}
	return domain.MenuItem{}, nil
}
