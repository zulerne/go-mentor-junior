package restaurant

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
	FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error)
}

type DeliveryProvider interface {
	Create(ctx context.Context, orderID uuid.UUID) error
	Start(ctx context.Context, orderID uuid.UUID) error
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

func (r *Service) GetOrders(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error) {
	orders, err := r.store.FindByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
