package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.orderStore.Find(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status == domain.Cancelled {
		return nil
	}
	if order.Status != domain.Pending {
		return domain.ErrInvalidOrderTransition
	}

	order.Status = domain.Cancelled

	_, err = s.orderStore.Update(ctx, order)
	return err
}
