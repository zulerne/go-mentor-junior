package restaurant_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestRejectOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)
	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrOrderNotFound)

	_, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
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

	assert.NoError(t, err)
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
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(nil)

	order, err := service.RejectOrder(context.Background(), testRestaurantID, testOrderID, "reason")

	assert.NoError(t, err)
	assert.Equal(t, domain.Rejected, order.Status)
	assert.Equal(t, "reason", order.RejectionReason)
	assert.Equal(t, testRestaurantID, order.RestaurantID)
}
