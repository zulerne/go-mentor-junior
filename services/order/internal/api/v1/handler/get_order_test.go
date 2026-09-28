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

func TestGetOrder_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+"invalid", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrder_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{"not found", domain.ErrNotFound, http.StatusNotFound},
		{"access denied", domain.ErrOrderAccessDenied, http.StatusForbidden},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			cust := NewMockCustomer(t)
			h := newHandler(cust, NewMockRestaurant(t))

			req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)
			req.Header.Set(customerIDHeader, testCustomerID.String())

			cust.EXPECT().GetOrder(mock.Anything, testOrderID).Return(domain.Order{}, tc.serviceErr)

			h.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func TestGetOrder_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetOrder(mock.Anything, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Pending,
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

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, handler.OrderResponse{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       string(domain.Pending),
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
