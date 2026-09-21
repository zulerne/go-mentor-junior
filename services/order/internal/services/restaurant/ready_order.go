package restaurant

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (r *Service) ReadyOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (domain.Order, error) {
	return domain.Order{}, nil

}
