package restaurant

import (
	"context"
	"log/slog"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id string) (domain.Order, error)
	FindByRestaurant(ctx context.Context, restaurantID string) ([]domain.Order, error)
}

type DeliveryProvider interface {
	Create(ctx context.Context, orderID string) error
	Start(ctx context.Context, orderID string) error
}

type Service struct {
	deliveryProvider DeliveryProvider
	store            OrderStore
	log              *slog.Logger
}

func New(store OrderStore, deliveryProvider DeliveryProvider, log *slog.Logger) *Service {
	log = log.With("component", "restaurant")
	r := &Service{
		store:            store,
		deliveryProvider: deliveryProvider,
		log:              log,
	}
	return r
}

func (r *Service) GetOrders(ctx context.Context, restaurantID string) ([]domain.Order, error) {
	orders, err := r.store.FindByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
