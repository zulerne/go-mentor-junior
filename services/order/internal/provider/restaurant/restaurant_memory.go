package restaurant

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type MemoryProvider struct {
	data map[uuid.UUID]map[uuid.UUID]domain.MenuItem
}

func NewMemoryProvider() *MemoryProvider {
	menuItems := map[uuid.UUID]map[uuid.UUID]domain.MenuItem{}
	menuItems[uuid.MustParse("restaurant-1")] = map[uuid.UUID]domain.MenuItem{
		uuid.MustParse("menu-item-1"): {
			ID:          uuid.MustParse("menu-item-1"),
			Name:        "Lasagna",
			Description: "Delicious lasagna pasta",
			PriceMinor:  100,
			Currency:    "USD",
			Available:   true,
		},
		uuid.MustParse("menu-item-2"): {
			ID:          uuid.MustParse("menu-item-2"),
			Name:        "Spaghetti",
			Description: "Delicious spaghetti pasta",
			PriceMinor:  80,
			Currency:    "USD",
			Available:   true,
		},
	}
	return &MemoryProvider{data: menuItems}
}

func (r *MemoryProvider) Find(_ context.Context, restaurantID uuid.UUID, menuItemID uuid.UUID) (domain.MenuItem, error) {
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
