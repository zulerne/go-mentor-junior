package customer

//go:generate go run github.com/vektra/mockery/v3@v3.8.0

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

// OrderStore is a persistent store for orders.
//
//mockery:generate: true
type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
	Update(ctx context.Context, order domain.Order) error
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
}

func New(orderStore OrderStore, cartStore CartStore, restaurantProvider RestaurantProvider, log *slog.Logger) *Service {
	log = log.With("component", "customer")
	c := &Service{
		orderStore:         orderStore,
		cartStore:          cartStore,
		restaurantProvider: restaurantProvider,
		log:                log,
	}

	return c
}
