package customer

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
	Update(ctx context.Context, order domain.Order) error
}

type CartStore interface {
	Find(ctx context.Context, customerID uuid.UUID) (domain.Cart, error)
	Update(ctx context.Context, customerID uuid.UUID, cart domain.Cart) error
}

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
