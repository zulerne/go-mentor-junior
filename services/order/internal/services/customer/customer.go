package customer

import (
	"context"
	"log/slog"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id string) (domain.Order, error)
}

type RestaurantProveder interface {
	Find(ctx context.Context, menuItemId string) (domain.MenuItem, error)
}

type Service struct {
	store              OrderStore
	restaurantProvider RestaurantProveder
	log                *slog.Logger
}

func New(store OrderStore, restaurantProvider RestaurantProveder, log *slog.Logger) *Service {
	log = log.With("component", "customer")
	c := &Service{
		store:              store,
		restaurantProvider: restaurantProvider,
		log:                log,
	}

	return c
}
