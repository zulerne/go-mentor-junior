package restaurant

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (r *Service) RejectOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, err := r.store.Find(ctx, orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if restaurantID != order.RestaurantID {
		return domain.Order{}, domain.ErrRestaurantIDMismatch
	}

	if order.Status == domain.Rejected {
		return order, nil
	}
	if order.Status != domain.Pending {
		return domain.Order{}, domain.ErrInvalidOrderStatus
	}

	order.Status = domain.Rejected
	order.RejectionReason = reason
	// TODO: remove when bd is connected
	order.UpdatedAt = time.Now().UTC()
	if err := r.store.Update(ctx, order); err != nil {
		return domain.Order{}, err
	}
	return order, nil
}
