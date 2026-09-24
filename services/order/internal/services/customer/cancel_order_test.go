package customer_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestCancelOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrOrderNotFound)

	_, err := service.CancelOrder(context.Background(), testCustomerID, testOrderID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrOrderNotFound, err)
}

func TestCancelOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	cStore := NewMockCartStore(t)
	rest := NewMockRestaurantProvider(t)

	service := newCustomer(oStore, cStore, rest)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:         testOrderID,
		CustomerID: uuid.New(),
		Status:     domain.Accepted}, nil)

	_, err := service.CancelOrder(context.Background(), testCustomerID, testOrderID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrOrderAccessDenied, err)
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

	order, err := service.CancelOrder(context.Background(), testCustomerID, testOrderID)

	assert.NoError(t, err)
	assert.Equal(t, domain.Cancelled, order.Status)
	assert.Equal(t, testOrderID, order.ID)
	assert.Equal(t, testCustomerID, order.CustomerID)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
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

	_, err := service.CancelOrder(context.Background(), testCustomerID, testOrderID)

	assert.Error(t, err)
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
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	order, err := service.CancelOrder(context.Background(), testCustomerID, testOrderID)

	assert.NoError(t, err)
	assert.Equal(t, domain.Cancelled, order.Status)
	assert.Equal(t, testOrderID, order.ID)
	assert.Equal(t, testCustomerID, order.CustomerID)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
}
