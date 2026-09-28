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

func TestGetAllOrders_MissingCustomerID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders", nil)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetAllOrders_Empty(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetAllOrders(mock.Anything, testCustomerID).Return([]domain.Order{}, nil)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	resp := decodeAllOrders(t, rec)
	assert.Empty(t, resp.Orders)
}

func TestGetAllOrders_Success(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetAllOrders(mock.Anything, testCustomerID).Return([]domain.Order{
		{
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
		},
	}, nil)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, handler.AllOrdersResponse{
		Orders: []handler.OrderResponse{{
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
		}},
	}, decodeAllOrders(t, rec))
}
