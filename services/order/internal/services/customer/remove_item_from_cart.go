package customer

import (
	"context"
	"slices"

	"github.com/google/uuid"
)

func (s *Service) RemoveItemFromCart(ctx context.Context, customerID uuid.UUID, menuItemID uuid.UUID) error {
	cart, err := s.cartStore.FindCart(ctx, customerID)
	if err != nil {
		return err
	}

	for i, cartItem := range cart.Items {
		if cartItem.MenuItemID == menuItemID {
			cart.Items = slices.Concat(cart.Items[:i], cart.Items[i+1:])
			cart.SubtotalMinor -= cartItem.UnitPriceMinor * int64(cartItem.Quantity)
			break
		}
	}

	if len(cart.Items) == 0 {
		cart.RestaurantID = uuid.Nil
		cart.Currency = ""
	}

	return s.cartStore.UpdateCart(ctx, customerID, cart)
}
