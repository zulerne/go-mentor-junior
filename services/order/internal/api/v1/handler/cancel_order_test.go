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

func TestCancelOrder_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+testOrderID.String()+"/cancel", nil)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCancelOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+"invalid"+"/cancel", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCancelOrder_NotFound(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+testOrderID.String()+"/cancel", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CancelOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{}, domain.ErrOrderNotFound)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCancelOrder_AccessDenied(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+testOrderID.String()+"/cancel", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CancelOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{}, domain.ErrOrderAccessDenied)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCancelOrder_InvalidTransition(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+testOrderID.String()+"/cancel", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CancelOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{}, domain.ErrInvalidOrderTransition)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestCancelOrder_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+testOrderID.String()+"/cancel", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CancelOrder(mock.Anything, testCustomerID, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Cancelled,
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
	assert.Equal(t, string(domain.Cancelled), resp.Status)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testMenuItemID.String(), resp.Items[0].MenuItemID)
	assert.Equal(t, "Test Item", resp.Items[0].Name)
	assert.EqualValues(t, 1, resp.Items[0].Quantity)
	assert.EqualValues(t, 100, resp.Items[0].UnitPriceMinor)
	assert.Equal(t, "", resp.Items[0].Instructions)
}
