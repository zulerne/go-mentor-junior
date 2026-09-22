package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error) {
	return nil, nil
}
