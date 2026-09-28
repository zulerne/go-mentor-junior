package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestAddItemToCart_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_InvalidMenuItemID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+"invalid", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_InvalidBody(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(
		context.Background(),
		"PUT",
		"/cart/items/"+testMenuItemID.String(),
		strings.NewReader("invalid"),
	)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_ValidationError_Quantity(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      11,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_ValidationError_Instructions(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	instructions := make([]byte, 251)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  string(instructions),
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{"cart limit exceeded", domain.ErrCartLimit, http.StatusUnprocessableEntity},
		{"restaurant conflict", domain.ErrNotFound, http.StatusConflict},
		{"item not available", domain.ErrMenuItemNotAvailable, http.StatusUnprocessableEntity},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			cust := NewMockCustomer(t)
			h := newHandler(cust, NewMockRestaurant(t))

			body := jsonBody(t, map[string]any{
				"restaurant_id": testRestaurantID,
				"quantity":      1,
				"instructions":  "",
			})
			req := httptest.NewRequestWithContext(
				context.Background(),
				"PUT",
				"/cart/items/"+testMenuItemID.String(),
				body,
			)
			req.Header.Set(customerIDHeader, testCustomerID.String())

			cust.EXPECT().AddItemToCart(mock.Anything, testCustomerID, testRestaurantID, domain.CartItem{
				MenuItemID:   testMenuItemID,
				Quantity:     1,
				Instructions: "",
			}).Return(domain.Cart{}, tc.serviceErr)

			h.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func TestAddItemToCart_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().AddItemToCart(mock.Anything, testCustomerID, testRestaurantID, domain.CartItem{
		MenuItemID:   testMenuItemID,
		Quantity:     1,
		Instructions: "",
	}).Return(domain.Cart{
		RestaurantID: testRestaurantID,
		Items: []domain.CartItem{
			{
				MenuItemID:     testMenuItemID,
				Name:           "Test Item",
				UnitPriceMinor: 0,
				Currency:       "USD",
				Quantity:       1,
				Instructions:   "",
			},
		},
		SubtotalMinor: 0,
		Currency:      "USD",
	}, nil)

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, handler.CartResponse{
		RestaurantID: new(testRestaurantID),
		Items: []handler.CartItem{{
			MenuItemID:     testMenuItemID,
			Name:           "Test Item",
			UnitPriceMinor: 0,
			Currency:       "USD",
			Quantity:       1,
			Instructions:   "",
		}},
		SubtotalMinor: 0,
		Currency:      "USD",
	}, decodeCart(t, rec))
}
