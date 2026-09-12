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
	FindCart(ctx context.Context, customerID string) (domain.Cart, error)
	AddItem(
		ctx context.Context,
		restaurantID, customerID string,
		menuItem domain.MenuItem,
		quantity int,
		instructions string,
	) error
}

type RestaurantProveder interface {
	Find(ctx context.Context, restaurantID, menuItemID string) (domain.MenuItem, error)
}

type Service struct {
	orderStore         OrderStore
	cartStore          CartStore
	restaurantProvider RestaurantProveder
	log                *slog.Logger
}

func New(orderStore OrderStore, cartStore CartStore, restaurantProvider RestaurantProveder, log *slog.Logger) *Service {
	log = log.With("component", "customer")
	c := &Service{
		orderStore:         orderStore,
		cartStore:          cartStore,
		restaurantProvider: restaurantProvider,
		log:                log,
	}

	return c
}

func (c *Service) GetCart(ctx context.Context, customerID string) (domain.Cart, error) {
	cart, err := c.cartStore.FindCart(ctx, customerID)
	if err != nil {
		return domain.Cart{}, err
	}
	return cart, nil
}
