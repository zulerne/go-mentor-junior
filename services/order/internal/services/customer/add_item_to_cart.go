package customer

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) AddItemToCart(ctx context.Context, customerID uuid.UUID, restaurantID uuid.UUID, menuItemID uuid.UUID, quantity int, instructions string) (domain.Cart, error) {
	cart, err := s.cartStore.FindCart(ctx, customerID)
	if err != nil || cart.RestaurantID == uuid.Nil {
		cart = domain.Cart{RestaurantID: restaurantID}
	} else if cart.RestaurantID != restaurantID {
		return domain.Cart{}, domain.ErrRestaurantIDMismatch
	}

	if len(instructions) > 250 {
		return domain.Cart{}, domain.ErrInvalidInstructions
	}

	totalQuantity := int32(quantity)
	for i, cartItem := range cart.Items {
		if cartItem.MenuItemID == menuItemID {
			cart.Items = slices.Concat(cart.Items[:i], cart.Items[i+1:])
			cart.SubtotalMinor -= cartItem.UnitPriceMinor * int64(cartItem.Quantity)
		} else {
			totalQuantity += cartItem.Quantity
		}
	}
	if len(cart.Items) >= 20 || totalQuantity > 50 {
		return domain.Cart{}, domain.ErrCartLimit
	}

	item, err := s.restaurantProvider.Find(ctx, restaurantID, menuItemID)
	if err != nil {
		return domain.Cart{}, err
	}

	if !item.Available {
		return domain.Cart{}, domain.ErrItemNotAvailable
	}

	cart.Items = append(cart.Items, domain.CartItem{
		MenuItemID:     menuItemID,
		Name:           item.Name,
		UnitPriceMinor: item.PriceMinor,
		Currency:       item.Currency,
		Quantity:       int32(quantity),
		Instructions:   instructions,
	})
	cart.Currency = item.Currency
	cart.SubtotalMinor += item.PriceMinor * int64(quantity)

	if err := s.cartStore.UpdateCart(ctx, customerID, cart); err != nil {
		return domain.Cart{}, err
	}

	return cart, nil
}
