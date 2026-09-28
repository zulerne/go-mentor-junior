package restaurant

//go:generate go run github.com/vektra/mockery/v3@v3.8.0

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

// OrderStore is a persistent store for orders.
//
//mockery:generate: true
type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
	FindByRestaurant(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error)
	Update(ctx context.Context, order domain.Order) (domain.Order, error)
}

// DeliveryProvider manages delivery lifecycle for orders.
//
//mockery:generate: true
type DeliveryProvider interface {
	Create(ctx context.Context, orderID uuid.UUID) error
	Start(ctx context.Context, orderID uuid.UUID) error
}

type Service struct {
	deliveryProvider DeliveryProvider
	store            OrderStore
	log              *slog.Logger
	now              func() time.Time
}

type Option func(*Service)

func WithNow(now func() time.Time) Option {
	return func(s *Service) { s.now = now }
}

func New(store OrderStore, deliveryProvider DeliveryProvider, log *slog.Logger, opts ...Option) *Service {
	log = log.With("component", "restaurant")
	s := &Service{
		store:            store,
		deliveryProvider: deliveryProvider,
		log:              log,
		now:              func() time.Time { return time.Now().UTC() },
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
