package restaurant

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
	FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error)
	Update(ctx context.Context, order domain.Order) error
}

type DeliveryProvider interface {
	Create(ctx context.Context, orderID uuid.UUID) error
	Start(ctx context.Context, orderID uuid.UUID) error
}

type Service struct {
	deliveryProvider DeliveryProvider
	store            OrderStore
	log              *slog.Logger
	mu               sync.Mutex
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
