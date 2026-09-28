package customer

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) AddItemToCart(
	ctx context.Context,
	customerID uuid.UUID,
	restaurantID uuid.UUID,
	data domain.CartItem,
) (domain.Cart, error) {
	cart, err := s.cartStore.Find(ctx, customerID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			return domain.Cart{}, err
		}
		cart = domain.Cart{RestaurantID: restaurantID}
	}
	if cart.RestaurantID == uuid.Nil {
		cart = domain.Cart{RestaurantID: restaurantID}
	} else if cart.RestaurantID != restaurantID {
		return domain.Cart{}, domain.ErrNotFound
	}

	itemIdx := -1
	totalQuantity := data.Quantity
	for i, cartItem := range cart.Items {
		if cartItem.MenuItemID == data.MenuItemID {
			itemIdx = i
		}
		totalQuantity += cartItem.Quantity
	}
	if len(cart.Items) >= 20 || totalQuantity > 50 {
		return domain.Cart{}, domain.ErrCartLimit
	}

	item, err := s.restaurantProvider.FindItem(ctx, cart.RestaurantID, data.MenuItemID)
	if err != nil {
		return domain.Cart{}, err
	}

	if !item.Available {
		return domain.Cart{}, domain.ErrMenuItemNotAvailable
	}

	if itemIdx != -1 {
		cart.Items[itemIdx].Name = item.Name
		cart.Items[itemIdx].UnitPriceMinor = item.PriceMinor
		cart.Items[itemIdx].Currency = item.Currency
		cart.Items[itemIdx].Quantity += data.Quantity
		cart.Items[itemIdx].Instructions = data.Instructions
	} else {
		cart.Items = append(cart.Items, domain.CartItem{
			MenuItemID:     data.MenuItemID,
			Name:           item.Name,
			UnitPriceMinor: item.PriceMinor,
			Currency:       item.Currency,
			Quantity:       data.Quantity,
			Instructions:   data.Instructions,
		})
	}
	if cart.Currency == "" {
		cart.Currency = item.Currency
	}
	cart.SubtotalMinor += item.PriceMinor * int64(data.Quantity)

	err = s.cartStore.Update(ctx, customerID, cart)
	if err != nil {
		return domain.Cart{}, err
	}

	return cart, nil
}
