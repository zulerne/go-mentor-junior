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

func TestRejectOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrNotFound)

	_, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestRejectOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: uuid.New(),
		Status:       domain.Pending,
	}, nil)

	_, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	assert.ErrorIs(t, err, domain.ErrOrderAccessDenied)
}

func TestRejectOrder_AlreadyRejected(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:              testOrderID,
		RestaurantID:    testRestaurantID,
		Status:          domain.Rejected,
		RejectionReason: "reason",
	}, nil)

	order, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	require.NoError(t, err)
	assert.Equal(t, domain.Rejected, order.Status)
}

func TestRejectOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		Status:       domain.Accepted,
	}, nil)

	_, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	assert.ErrorIs(t, err, domain.ErrInvalidOrderTransition)
}

func TestRejectOrder_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Pending,
	}, nil)
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.Order{
		ID:              testOrderID,
		RestaurantID:    testRestaurantID,
		CustomerID:      testCustomerID,
		Status:          domain.Rejected,
		RejectionReason: "reason",
		UpdatedAt:       fixedNow,
	}, nil)

	order, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	require.NoError(t, err)
	assert.Equal(t, domain.Order{
		ID:              testOrderID,
		RestaurantID:    testRestaurantID,
		CustomerID:      testCustomerID,
		Status:          domain.Rejected,
		RejectionReason: "reason",
		UpdatedAt:       fixedNow,
	}, order)
}
