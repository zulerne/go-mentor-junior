package restaurant

import (
	"context"
	"log/slog"

	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type OrderStore interface {
	Find(ctx context.Context, id string) (domain.Order, error)
}

type DeliveryProvider interface {
	Create(ctx context.Context, order_id string) error
	Start(ctx context.Context, order_id string) error
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
