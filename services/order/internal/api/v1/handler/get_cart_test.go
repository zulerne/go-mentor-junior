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

func TestGetCart_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/cart", nil)
	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetCart_EmptyCart(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/cart", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetCart(mock.Anything, testCustomerID).Return(domain.Cart{}, domain.ErrCartNotFound)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeCart(t, rec)
	assert.Empty(t, resp.Items)
	assert.EqualValues(t, 0, resp.SubtotalMinor)
}

func TestGetCart_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/cart", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetCart(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 150,
		Currency:      "USD",
	}, nil)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeCart(t, rec)
	assert.Equal(t, testRestaurantID.String(), resp.RestaurantID)
	assert.Equal(t, "USD", resp.Currency)
	assert.EqualValues(t, 150, resp.SubtotalMinor)
	assert.Empty(t, resp.Items)
}
