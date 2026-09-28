package restaurant_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestReadyOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrNotFound)

	_, err := service.ReadyOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestReadyOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: uuid.New(),
		Status:       domain.Preparing,
	}, nil)

	_, err := service.ReadyOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrOrderAccessDenied)
}

func TestReadyOrder_AlreadyReadyForPickup(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		Status:       domain.ReadyForPickup,
	}, nil)

	order, err := service.ReadyOrder(context.Background(), testRestaurantID, testOrderID)

	require.NoError(t, err)
	assert.Equal(t, domain.ReadyForPickup, order.Status)
}

func TestReadyOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		Status:       domain.Accepted,
	}, nil)

	_, err := service.ReadyOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrInvalidOrderTransition)
}

func TestReadyOrder_DeliveryTransitionFails(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		Status:       domain.Preparing,
	}, nil)
	delivery.EXPECT().Start(mock.Anything, testOrderID).Return(domain.ErrDeliveryProvider)

	_, err := service.ReadyOrder(context.Background(), testRestaurantID, testOrderID)

	assert.ErrorIs(t, err, domain.ErrDeliveryProvider)
}

func TestReadyOrder_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Preparing,
	}, nil)
	delivery.EXPECT().Start(mock.Anything, testOrderID).Return(nil)
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.Order{
		ID:             testOrderID,
		RestaurantID:   testRestaurantID,
		CustomerID:     testCustomerID,
		Status:         domain.ReadyForPickup,
		DeliveryStatus: domain.InDelivery,
		UpdatedAt:      fixedNow,
	}, nil)

	order, err := service.ReadyOrder(context.Background(), testRestaurantID, testOrderID)

	require.NoError(t, err)
	assert.Equal(t, domain.Order{
		ID:             testOrderID,
		RestaurantID:   testRestaurantID,
		CustomerID:     testCustomerID,
		Status:         domain.ReadyForPickup,
		DeliveryStatus: domain.InDelivery,
		UpdatedAt:      fixedNow,
	}, order)
}
