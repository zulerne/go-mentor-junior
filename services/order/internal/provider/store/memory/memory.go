package memory

import (
	"context"
	"time"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore struct {
	data map[string]domain.Order
}

func NewOrderStore() *OrderStore {
	data := make(map[string]domain.Order)

	data["test"] = domain.Order{
		ID:              "test",
		CustomerID:      "",
		RestaurantID:    "",
		Status:          domain.Pending,
		Items:           nil,
		SubtotalMinor:   0,
		Currency:        "",
		DeliveryAddress: "",
		RejectionReason: nil,
		DeliveryStatus:  nil,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	return &OrderStore{
		data: data,
	}
}

func (s *OrderStore) Find(ctx context.Context, id string) (domain.Order, error) {
	order, ok := s.data[id]
	if !ok {
		return domain.Order{}, nil
	}
	return order, nil
}
