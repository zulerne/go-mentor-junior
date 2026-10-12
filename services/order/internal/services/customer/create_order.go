package customer

import (
	"context"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (s *Service) CreateOrder(ctx context.Context, customerID uuid.UUID, deliveryAddress string) (domain.Order, error) {
	var result domain.Order
	err := s.withinCart(ctx, customerID, func(carts CartStore, orders OrderStore) error {
		var createErr error
		result, createErr = s.createOrder(ctx, customerID, deliveryAddress, carts, orders)
		return createErr
	})
	if err != nil {
		return domain.Order{}, err
	}
	return result, nil
}

func (s *Service) createOrder(ctx context.Context, customerID uuid.UUID, deliveryAddress string,
	carts CartStore, orders OrderStore) (domain.Order, error) {
	cart, err := carts.Find(ctx, customerID)
	if err != nil {
		return domain.Order{}, err
	}
	if len(cart.Items) == 0 {
		return domain.Order{}, domain.ErrCartEmpty
	}

	restaurant, err := s.restaurantProvider.Find(ctx, cart.RestaurantID)
	if err != nil {
		return domain.Order{}, err
	}
	if !restaurant.AcceptingOrders {
		return domain.Order{}, domain.ErrRestaurantNotAcceptingOrders
	}

	cart.SubtotalMinor = 0
	var item domain.MenuItem
	for i, cartItem := range cart.Items {
		item, err = s.restaurantProvider.FindItem(ctx, restaurant.ID, cartItem.MenuItemID)
		if err != nil {
			return domain.Order{}, err
		}
		if !item.Available {
			return domain.Order{}, domain.ErrMenuItemNotAvailable
		}
		cart.Items[i].UnitPriceMinor = item.PriceMinor
		cart.SubtotalMinor += item.PriceMinor * int64(cartItem.Quantity)
	}

	if cart.SubtotalMinor < restaurant.MinimumOrderMinor {
		return domain.Order{}, domain.ErrMinOrderNotReached
	}

	orderItems := make([]domain.OrderItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		orderItems = append(orderItems, domain.OrderItem{
			MenuItemID:     item.MenuItemID,
			Name:           item.Name,
			UnitPriceMinor: item.UnitPriceMinor,
			Quantity:       item.Quantity,
			Instructions:   item.Instructions,
		})
	}

	now := s.now()

	order := domain.Order{
		ID:              s.newID(),
		CustomerID:      customerID,
		RestaurantID:    cart.RestaurantID,
		Status:          domain.Pending,
		Items:           orderItems,
		SubtotalMinor:   cart.SubtotalMinor,
		Currency:        cart.Currency,
		DeliveryAddress: deliveryAddress,
		RejectionReason: "",
		DeliveryStatus:  domain.UnspecifiedOrderDeliveryStatus,
		CreatedAt:       now,
	}

	order, err = orders.Update(ctx, order)
	if err != nil {
		return domain.Order{}, err
	}

	err = carts.Update(ctx, customerID, domain.Cart{})
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
