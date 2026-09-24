package customer_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestRemoveItemFromCart_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Quantity:       1,
				UnitPriceMinor: 100,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)
	cStore.EXPECT().Update(mock.Anything, testCustomerID, domain.Cart{}).Return(nil)

	err := service.RemoveItemFromCart(context.Background(), testCustomerID, testMenuItemID)

	assert.NoError(t, err)
}

func TestRemoveItemFromCart_ItemNotInCart(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	cStore.EXPECT().Find(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Quantity:       1,
				UnitPriceMinor: 100,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)
	cStore.EXPECT().Update(mock.Anything, testCustomerID, mock.Anything).Return(nil)

	err := service.RemoveItemFromCart(context.Background(), testCustomerID, uuid.New())

	assert.NoError(t, err)
}
