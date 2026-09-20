package store

import (
	"context"
	"time"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type MemoryOrderStore struct {
	data map[string]domain.Order
}

func NewMemoryOrderStore() *MemoryOrderStore {
	data := make(map[string]domain.Order)

	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC)

	data["test"] = domain.Order{
		ID:              "test",
		CustomerID:      "",
		RestaurantID:    "1",
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

func (s *MemoryOrderStore) Find(_ context.Context, id string) (domain.Order, error) {
	order, ok := s.data[id]
	if !ok {
		return domain.Order{}, nil
	}
	return order, nil
}

func (s *MemoryOrderStore) FindByRestaurant(_ context.Context, restaurantID string) ([]domain.Order, error) {
	var orders []domain.Order
	for _, order := range s.data {
		if order.RestaurantID == restaurantID {
			orders = append(orders, order)
		}
	}
	return orders, nil
}
