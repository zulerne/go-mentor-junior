package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) GetOrder(ctx context.Context, customerID uuid.UUID, orderID uuid.UUID) (domain.Order, error) {
	order, err := s.orderStore.Find(ctx, orderID)
	if err != nil {
		return domain.Order{}, err
	}
	if order.CustomerID != customerID {
		return domain.Order{}, domain.ErrOrderAccessDenied
	}
	return order, nil
}
