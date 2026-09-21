package store

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type MemoryOrderStore struct {
	data map[uuid.UUID]domain.Order
}

func NewMemoryOrderStore() *MemoryOrderStore {
	data := make(map[uuid.UUID]domain.Order)

	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC)
	// restId: fe70cb38-d10d-452c-8860-1af5936f7037
	// id1: 58185771-1f9f-4d71-9609-bdacea1deb2e
	// id2: 4c5344f5-33d7-4c05-bbbc-a0edc4178e7e
	// userId: ef55f77a-7738-426f-93a4-f78f2baf7970
	// orderId: 2f5efdc6-c936-4fdf-a681-1e21e73e6e71

	orderId := uuid.MustParse("2f5efdc6-c936-4fdf-a681-1e21e73e6e71")
	customerId := uuid.MustParse("ef55f77a-7738-426f-93a4-f78f2baf7970")
	restaurantId := uuid.MustParse("fe70cb38-d10d-452c-8860-1af5936f7037")
	data[orderId] = domain.Order{
		ID:              orderId,
		CustomerID:      customerId,
		RestaurantID:    restaurantId,
		Status:          domain.Pending,
		Items:           nil,
		SubtotalMinor:   0,
		Currency:        "",
		DeliveryAddress: "",
		RejectionReason: "",
		DeliveryStatus:  domain.UnspecifiedOrderDeliveryStatus,
		CreatedAt:       date,
		UpdatedAt:       date,
	}
	return &MemoryOrderStore{
		data: data,
	}
}

func (s *MemoryOrderStore) Find(_ context.Context, id uuid.UUID) (domain.Order, error) {
	order, ok := s.data[id]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}
	return order, nil
}

func (s *MemoryOrderStore) FindByRestaurant(_ context.Context, restaurantID uuid.UUID) ([]domain.Order, error) {
	var orders []domain.Order
	for _, order := range s.data {
		if order.RestaurantID == restaurantID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}

func (s *MemoryOrderStore) Update(_ context.Context, order domain.Order) error {
	s.data[order.ID] = order
	return nil
}
