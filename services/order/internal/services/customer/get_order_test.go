package customer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestGetOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrNotFound)

	_, err := service.GetOrder(context.Background(), testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestGetOrder_Success(t *testing.T) {
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
		Currency:     "USD",
	}, nil)

	order, err := service.GetOrder(context.Background(), testOrderID)

	require.NoError(t, err)
	assert.Equal(t, domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Pending,
		Currency:     "USD",
	}, order)
}
