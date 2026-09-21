package customer

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id uuid.UUID) (domain.Order, error)
}

type CartStore interface {
	FindCart(ctx context.Context, customerID uuid.UUID) (domain.Cart, error)
	AddItem(
		ctx context.Context,
		restaurantID uuid.UUID,
		customerID uuid.UUID,
		menuItem domain.MenuItem,
		quantity int,
		instructions string,
	) error
}

type RestaurantProvider interface {
	Find(ctx context.Context, restaurantID uuid.UUID, menuItemID uuid.UUID) (domain.MenuItem, error)
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
