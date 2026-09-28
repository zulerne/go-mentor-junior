package customer

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
	Update(ctx context.Context, order domain.Order) (domain.Order, error)
	GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error)
}

// CartStore is a persistent store for carts.
//
//mockery:generate: true
type CartStore interface {
	Find(ctx context.Context, customerID uuid.UUID) (domain.Cart, error)
	Update(ctx context.Context, customerID uuid.UUID, cart domain.Cart) error
}

// RestaurantProvider provides access to restaurant and menu data.
//
//mockery:generate: true
type RestaurantProvider interface {
	Find(ctx context.Context, restaurantID uuid.UUID) (domain.Restaurant, error)
	FindItem(ctx context.Context, restaurantID uuid.UUID, menuItemID uuid.UUID) (domain.MenuItem, error)
}

type Service struct {
	orderStore         OrderStore
	cartStore          CartStore
	restaurantProvider RestaurantProvider
	log                *slog.Logger
	now                func() time.Time
	newID              func() uuid.UUID
}

type Option func(*Service)

func WithNow(now func() time.Time) Option {
	return func(s *Service) { s.now = now }
}

func WithNewID(newID func() uuid.UUID) Option {
	return func(s *Service) { s.newID = newID }
}

func New(
	orderStore OrderStore,
	cartStore CartStore,
	restaurantProvider RestaurantProvider,
	log *slog.Logger,
	opts ...Option,
) *Service {
	log = log.With("component", "customer")
	s := &Service{
		orderStore:         orderStore,
		cartStore:          cartStore,
		restaurantProvider: restaurantProvider,
		log:                log,
		now:                func() time.Time { return time.Now().UTC() },
		newID:              uuid.New,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
