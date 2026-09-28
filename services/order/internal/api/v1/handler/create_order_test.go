package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestCreateOrder_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", strings.NewReader("text"))

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_InvalidBody(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", strings.NewReader("text"))
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_InvalidDeliveryAddress(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "    ",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{"empty cart", domain.ErrCartEmpty, http.StatusUnprocessableEntity},
		{"restaurant not found", domain.ErrNotFound, http.StatusNotFound},
		{"restaurant not accepting", domain.ErrRestaurantNotAcceptingOrders, http.StatusBadRequest},
		{"minimum order not reached", domain.ErrMinOrderNotReached, http.StatusBadRequest},
		{"item not available", domain.ErrMenuItemNotAvailable, http.StatusUnprocessableEntity},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			cust := NewMockCustomer(t)
			h := newHandler(cust, NewMockRestaurant(t))

			body := jsonBody(t, map[string]any{"delivery_address": "address"})
			req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
			req.Header.Set(customerIDHeader, testCustomerID.String())

			cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{}, tc.serviceErr)

			h.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func TestCreateOrder_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "address",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{
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

	assert.Equal(t, http.StatusCreated, rec.Code)
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
