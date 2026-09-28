package customer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestCancelOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrNotFound)

	err := service.CancelOrder(context.Background(), testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestCancelOrder_AlreadyCancelled(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Cancelled,
	}, nil)

	err := service.CancelOrder(context.Background(), testOrderID)

	require.NoError(t, err)
}

func TestCancelOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Accepted,
	}, nil)

	err := service.CancelOrder(context.Background(), testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrInvalidOrderTransition, err)
}

func TestCancelOrder_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Pending,
	}, nil)
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.Order{}, nil)

	err := service.CancelOrder(context.Background(), testOrderID)

	require.NoError(t, err)
}
