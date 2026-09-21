package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (c *Service) GetCart(ctx context.Context, customerID uuid.UUID) (domain.Cart, error) {
	cart, err := c.cartStore.FindCart(ctx, customerID)
	if err != nil {
		return domain.Cart{}, err
	}
	return cart, nil
}
