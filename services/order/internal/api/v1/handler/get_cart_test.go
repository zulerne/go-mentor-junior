package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func TestGetCart_MissingCustomerID(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/cart", nil)
	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestGetCart_EmptyCart(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/cart", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetCart(mock.Anything, testCustomerID).Return(domain.Cart{}, domain.ErrNotFound)

	h.Routes().ServeHTTP(rec, req)

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
	h := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "GET", "/cart", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().GetCart(mock.Anything, testCustomerID).Return(domain.Cart{
		RestaurantID:  testRestaurantID,
		Items:         []domain.CartItem{},
		SubtotalMinor: 150,
		Currency:      "USD",
	}, nil)

	h.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, handler.CartResponse{
		RestaurantID:  new(testRestaurantID),
		Items:         []handler.CartItem{},
		SubtotalMinor: 150,
		Currency:      "USD",
	}, decodeCart(t, rec))
}
