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

func TestAddItemToCart_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_InvalidMenuItemID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+"invalid", body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_InvalidBody(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), strings.NewReader("inmvalid"))
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_ValidationError_Quantity(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      11,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_ValidationError_Instructions(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	instructions := make([]byte, 251)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  string(instructions),
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAddItemToCart_CartLimitExceeded(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().AddItemToCart(mock.Anything, testCustomerID, testRestaurantID, testMenuItemID, int32(1), "").Return(domain.Cart{}, domain.ErrCartLimit)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestAddItemToCart_RestaurantConflict(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().AddItemToCart(mock.Anything, testCustomerID, testRestaurantID, testMenuItemID, int32(1), "").Return(domain.Cart{}, domain.ErrOrderAccessDenied)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAddItemToCart_ItemNotAvailable(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().AddItemToCart(mock.Anything, testCustomerID, testRestaurantID, testMenuItemID, int32(1), "").Return(domain.Cart{}, domain.ErrMenuItemNotAvailable)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestAddItemToCart_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	body := jsonBody(t, map[string]any{
		"restaurant_id": testRestaurantID,
		"quantity":      1,
		"instructions":  "",
	})
	req := httptest.NewRequestWithContext(context.Background(), "PUT", "/cart/items/"+testMenuItemID.String(), body)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().AddItemToCart(mock.Anything, testCustomerID, testRestaurantID, testMenuItemID, int32(1), "").Return(domain.Cart{
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

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeCart(t, rec)
	assert.Equal(t, testRestaurantID.String(), resp.RestaurantID)
	assert.Equal(t, "USD", resp.Currency)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testMenuItemID.String(), resp.Items[0].MenuItemID)
	assert.Equal(t, "Test Item", resp.Items[0].Name)
	assert.EqualValues(t, 1, resp.Items[0].Quantity)
}
