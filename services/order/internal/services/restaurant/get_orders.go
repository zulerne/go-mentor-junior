package restaurant

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (r *Service) GetOrders(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error) {
	orders, err := r.store.FindByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, domain.ErrRestaurantNotFound
	}
	return orders, nil
}
