package customer

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) CancelOrder(ctx context.Context, customerID, orderID uuid.UUID) (domain.Order, error) {
	order, err := s.orderStore.Find(ctx, orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if order.CustomerID != customerID {
		return domain.Order{}, domain.ErrOrderAccessDenied
	}
	if order.Status == domain.Cancelled {
		return order, nil
	}
	if order.Status != domain.Pending {
		return domain.Order{}, domain.ErrInvalidOrderStatus
	}

	order.Status = domain.Cancelled
	order.UpdatedAt = time.Now()

	err = s.orderStore.Update(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}
