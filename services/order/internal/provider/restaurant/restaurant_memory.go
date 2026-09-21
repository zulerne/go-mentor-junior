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
	// restaurantId: fe70cb38-d10d-452c-8860-1af5936f7037
	// id1: 58185771-1f9f-4d71-9609-bdacea1deb2e
	// id2: 4c5344f5-33d7-4c05-bbbc-a0edc4178e7e
	// userId: ef55f77a-7738-426f-93a4-f78f2baf7970
	// orderId: of5efdc6-c936-4fdf-a681-1e21e73e6e71
	restaurantId := uuid.MustParse("fe70cb38-d10d-452c-8860-1af5936f7037")
	id1 := uuid.MustParse("58185771-1f9f-4d71-9609-bdacea1deb2e")
	id2 := uuid.MustParse("4c5344f5-33d7-4c05-bbbc-a0edc4178e7e")
	menuItems[restaurantId] = map[uuid.UUID]domain.MenuItem{
		id1: {
			ID:          id1,
			Name:        "Lasagna",
			Description: "Delicious lasagna pasta",
			PriceMinor:  100,
			Currency:    "USD",
			Available:   true,
		},
		id2: {
			ID:          id2,
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
