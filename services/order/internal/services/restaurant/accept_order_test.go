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

func TestAcceptOrder_NotFound(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)

	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{}, domain.ErrNotFound)

	_, err := service.AcceptOrder(context.Background(), testRestaurantID, testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrNotFound, err)
}

func TestAcceptOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)

	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: uuid.New(),
		CustomerID:   testCustomerID,
		Status:       domain.Rejected,
		Items: []domain.OrderItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	_, err := service.AcceptOrder(context.Background(), testRestaurantID, testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrOrderAccessDenied, err)
}

func TestAcceptOrder_AlreadyAccepted(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)

	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Accepted,
		Items: []domain.OrderItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	_, err := service.AcceptOrder(context.Background(), testRestaurantID, testOrderID)

	require.NoError(t, err)
}

func TestAcceptOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)

	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Rejected,
		Items: []domain.OrderItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	_, err := service.AcceptOrder(context.Background(), testRestaurantID, testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrInvalidOrderTransition, err)
}

func TestAcceptOrder_DeliveryCreationFails(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)

	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Pending,
		Items: []domain.OrderItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	delivery.EXPECT().Create(mock.Anything, testOrderID).Return(domain.ErrDeliveryProvider)

	_, err := service.AcceptOrder(context.Background(), testRestaurantID, testOrderID)

	require.Error(t, err)
	assert.Equal(t, domain.ErrDeliveryProvider, err)
}

func TestAcceptOrder_Success(t *testing.T) {
	t.Parallel()

	oStore := NewMockOrderStore(t)
	delivery := NewMockDeliveryProvider(t)

	service := newRestaurant(oStore, delivery)

	oStore.EXPECT().Find(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		RestaurantID: testRestaurantID,
		CustomerID:   testCustomerID,
		Status:       domain.Pending,
		Items: []domain.OrderItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "item",
				UnitPriceMinor: 100,
				Quantity:       1,
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	delivery.EXPECT().Create(mock.Anything, testOrderID).Return(nil)
	oStore.EXPECT().Update(mock.Anything, mock.Anything).Return(domain.Order{
		ID:             testOrderID,
		RestaurantID:   testRestaurantID,
		CustomerID:     testCustomerID,
		Status:         domain.Accepted,
		DeliveryStatus: domain.WaitingForPreparation,
		Items: []domain.OrderItem{{
			MenuItemID:     testMenuItemID,
			Name:           "item",
			UnitPriceMinor: 100,
			Quantity:       1,
		}},
		SubtotalMinor: 100,
		Currency:      "USD",
		UpdatedAt:     fixedNow,
	}, nil)

	order, err := service.AcceptOrder(context.Background(), testRestaurantID, testOrderID)

	require.NoError(t, err)
	assert.Equal(t, domain.Order{
		ID:             testOrderID,
		RestaurantID:   testRestaurantID,
		CustomerID:     testCustomerID,
		Status:         domain.Accepted,
		DeliveryStatus: domain.WaitingForPreparation,
		Items: []domain.OrderItem{{
			MenuItemID:     testMenuItemID,
			Name:           "item",
			UnitPriceMinor: 100,
			Quantity:       1,
		}},
		SubtotalMinor: 100,
		Currency:      "USD",
		UpdatedAt:     fixedNow,
	}, order)
}
