package restaurant

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (r *Service) ReadyOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (domain.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, err := r.store.Find(ctx, orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if restaurantID != order.RestaurantID {
		return domain.Order{}, domain.ErrOrderAccessDenied
	}

	if order.Status == domain.ReadyForPickup {
		return order, nil
	}
	if order.Status != domain.Preparing {
		return domain.Order{}, domain.ErrInvalidOrderTransition
	}

	err = r.deliveryProvider.Start(ctx, orderID)
	if err != nil {
		return domain.Order{}, domain.ErrDeliveryProvider
	}

	order.Status = domain.ReadyForPickup
	order.DeliveryStatus = domain.InDelivery
	// TODO: remove when bd is connected
	order.UpdatedAt = time.Now().UTC()
	err = r.store.Update(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}
