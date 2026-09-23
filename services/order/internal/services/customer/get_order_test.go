package customer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestGetOrder_NotFound(t *testing.T) {
	t.Parallel()
}

func TestGetOrder_AccessDenied(t *testing.T) {
	t.Parallel()
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

	order, err := service.GetOrder(context.Background(), testCustomerID, testOrderID)

	assert.NoError(t, err)
	assert.Equal(t, testOrderID, order.ID)
	assert.Equal(t, domain.Pending, order.Status)
	assert.Equal(t, testCustomerID, order.CustomerID)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
}
