package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) GetCart(ctx context.Context, customerID uuid.UUID) (domain.Cart, error) {
	cart, err := s.cartStore.Find(ctx, customerID)
	if err != nil {
		return domain.Cart{}, err
	}
	return cart, nil
}
