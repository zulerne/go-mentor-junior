package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) CancelOrder(ctx context.Context, customerID, orderID uuid.UUID) (domain.Order, error) {
	return domain.Order{}, nil
}
