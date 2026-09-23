package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	assert.Len(t, resp.Orders, 0)
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

	resp := decodeAllOrders(t, rec)
	assert.Equal(t, testOrderID.String(), resp.Orders[0].ID)
	assert.Equal(t, testCustomerID.String(), resp.Orders[0].CustomerID)
	assert.Equal(t, testRestaurantID.String(), resp.Orders[0].RestaurantID)
	assert.Equal(t, string(domain.Pending), resp.Orders[0].Status)
	assert.Len(t, resp.Orders[0].Items, 1)
	assert.Equal(t, testMenuItemID.String(), resp.Orders[0].Items[0].MenuItemID)
	assert.Equal(t, "Test Item", resp.Orders[0].Items[0].Name)
	assert.EqualValues(t, 1, resp.Orders[0].Items[0].Quantity)
	assert.EqualValues(t, 100, resp.Orders[0].Items[0].UnitPriceMinor)
	assert.Equal(t, "", resp.Orders[0].Items[0].Instructions)
}
