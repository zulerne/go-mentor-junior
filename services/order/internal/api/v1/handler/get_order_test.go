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

func TestGetOrder_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+"invalid", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetOrder_NotFound(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{}, domain.ErrOrderNotFound)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetOrder_AccessDenied(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{}, domain.ErrOrderAccessDenied)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetOrder_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/orders/"+testOrderID.String(), nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{
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

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeOrder(t, rec)
	assert.Equal(t, testOrderID.String(), resp.ID)
	assert.Equal(t, testCustomerID.String(), resp.CustomerID)
	assert.Equal(t, testRestaurantID.String(), resp.RestaurantID)
	assert.Equal(t, string(domain.Pending), resp.Status)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testMenuItemID.String(), resp.Items[0].MenuItemID)
	assert.Equal(t, "Test Item", resp.Items[0].Name)
	assert.EqualValues(t, 1, resp.Items[0].Quantity)
	assert.EqualValues(t, 100, resp.Items[0].UnitPriceMinor)
	assert.Equal(t, "", resp.Items[0].Instructions)
}
