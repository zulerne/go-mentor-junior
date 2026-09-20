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

	data[uuid.MustParse("test")] = domain.Order{
		ID:              uuid.MustParse("test"),
		CustomerID:      uuid.MustParse(""),
		RestaurantID:    uuid.MustParse("1"),
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
		return domain.Order{}, nil
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
