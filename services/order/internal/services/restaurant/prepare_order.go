package restaurant

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (r *Service) PrepareOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (domain.Order, error) {
	order, err := r.store.Find(ctx, orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if restaurantID != order.RestaurantID {
		return domain.Order{}, domain.ErrOrderAccessDenied
	}

	if order.Status == domain.Preparing {
		return order, nil
	}
	if order.Status != domain.Accepted {
		return domain.Order{}, domain.ErrInvalidOrderTransition
	}

	order.Status = domain.Preparing
	updated, err := r.store.Update(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}
	return updated, nil
}
