package customer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestCancelOrder_NotFound(t *testing.T) {
	t.Parallel()
}

func TestCancelOrder_AccessDenied(t *testing.T) {
	t.Parallel()
}

func TestCancelOrder_AlreadyCancelled(t *testing.T) {
	t.Parallel()
}

func TestCancelOrder_WrongStatus(t *testing.T) {
	t.Parallel()
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
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	order, err := service.CancelOrder(context.Background(), testCustomerID, testOrderID)

	assert.NoError(t, err)
	assert.Equal(t, domain.Cancelled, order.Status)
	assert.Equal(t, testOrderID, order.ID)
	assert.Equal(t, testCustomerID, order.CustomerID)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
}
