package restaurant_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestPrepareOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrOrderNotFound)

	_, err := service.PrepareOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
}

func TestPrepareOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: uuid.New(),
		Status:       domain.Accepted,
	}, nil)

	_, err := service.PrepareOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrOrderAccessDenied)
}

func TestPrepareOrder_AlreadyPreparing(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		Status:       domain.Preparing,
	}, nil)

	order, err := service.PrepareOrder(context.Background(), testRestaurantID, testOrderID)

	assert.NoError(t, err)
	assert.Equal(t, domain.Preparing, order.Status)
}

func TestPrepareOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		Status:       domain.Pending,
	}, nil)

	_, err := service.PrepareOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrInvalidOrderTransition)
}

func TestPrepareOrder_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Accepted,
	}, nil)
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	order, err := service.PrepareOrder(context.Background(), testRestaurantID, testOrderID)

	assert.NoError(t, err)
	assert.Equal(t, domain.Preparing, order.Status)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
}
