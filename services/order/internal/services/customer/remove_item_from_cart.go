package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) RemoveItemFromCart(ctx context.Context, customerID uuid.UUID, itemID uuid.UUID) (domain.Cart, error) {
	return domain.Cart{}, nil
}
