package customer

import (
	"context"
	"log/slog"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id string) (domain.Order, error)
}

type CartStore interface {
	FindCart(ctx context.Context, customerId string) (domain.Cart, error)
	AddItem(ctx context.Context, restaurantId, customerId string, menuItem domain.MenuItem, quantity int, instructions string) error
}

type RestaurantProveder interface {
	Find(ctx context.Context, restaurantId, menuItemId string) (domain.MenuItem, error)
}

type Service struct {
	OrderStore         OrderStore
	CartStore          CartStore
	RestaurantProvider RestaurantProveder
	log                *slog.Logger
}

func New(orderStore OrderStore, cartStore CartStore, restaurantProvider RestaurantProveder, log *slog.Logger) *Service {
	log = log.With("component", "customer")
	c := &Service{
		OrderStore:         orderStore,
		CartStore:          cartStore,
		RestaurantProvider: restaurantProvider,
		log:                log,
	}

	return c
}
