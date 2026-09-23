package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestCreateOrder_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", strings.NewReader("text"))

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_InvalidBody(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", strings.NewReader("text"))
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_EmptyCart(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "address",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{}, domain.ErrCartEmpty)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateOrder_RestaurantNotFound(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "address",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{}, domain.ErrRestaurantNotFound)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateOrder_RestaurantNotAccepting(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "address",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{}, domain.ErrRestaurantNotAcceptingOrders)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_MinimumOrderNotReached(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "address",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{}, domain.ErrMinOrderNotReached)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_ItemNotAvailable(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "address",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CreateOrder(mock.Anything, testCustomerID, "address").Return(domain.Order{}, domain.ErrMenuItemNotAvailable)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateOrder_InvalidDeliveryAddress(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"delivery_address": "    ",
	})
	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOrder_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

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

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
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
