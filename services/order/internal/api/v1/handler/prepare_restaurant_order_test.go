package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestPrepareRestaurantOrder_MissingRestaurantID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(
		context.Background(),
		"POST",
		"/restaurant/orders/"+testOrderID.String()+"/start-preparation",
		nil,
	)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPrepareRestaurantOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(
		context.Background(),
		"POST",
		"/restaurant/orders/"+"invalid"+"/start-preparation",
		nil,
	)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPrepareRestaurantOrder_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{"not found", domain.ErrNotFound, http.StatusNotFound},
		{"access denied", domain.ErrOrderAccessDenied, http.StatusForbidden},
		{"invalid transition", domain.ErrInvalidOrderTransition, http.StatusConflict},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			rest := NewMockRestaurant(t)
			hand := newHandler(NewMockCustomer(t), rest)

			req := httptest.NewRequestWithContext(
				context.Background(),
				"POST",
				"/restaurant/orders/"+testOrderID.String()+"/start-preparation",
				nil,
			)
			req.Header.Set(restaurantIDHeader, testRestaurantID.String())

			rest.EXPECT().
				PrepareOrder(mock.Anything, testRestaurantID, testOrderID).
				Return(domain.Order{}, tc.serviceErr)

			hand.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func TestPrepareRestaurantOrder_Success(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(
		context.Background(),
		"POST",
		"/restaurant/orders/"+testOrderID.String()+"/start-preparation",
		nil,
	)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().PrepareOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Preparing,
		Items: []domain.OrderItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "Test Item",
				UnitPriceMinor: 100,
				Quantity:       1,
				Instructions:   "",
			},
		},
		SubtotalMinor: 100,
		Currency:      "USD",
	}, nil)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, handler.OrderResponse{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       string(domain.Preparing),
		Items: []handler.OrderItem{{
			MenuItemID:     testMenuItemID,
			Name:           "Test Item",
			UnitPriceMinor: 100,
			Quantity:       1,
			Instructions:   "",
		}},
		SubtotalMinor:   100,
		Currency:        "USD",
		DeliveryAddress: "",
		RejectionReason: "",
		DeliveryStatus:  "",
		CreatedAt:       time.Time{}.UTC(),
		UpdatedAt:       time.Time{}.UTC(),
	}, decodeOrder(t, rec))
}
