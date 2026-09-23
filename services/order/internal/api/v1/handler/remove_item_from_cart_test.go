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

func TestRemoveItemFromCart_MissingCustomerID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "DELETE", "/cart/items/"+testMenuItemID.String(), nil)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRemoveItemFromCart_InvalidMenuItemID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "DELETE", "/cart/items/invalid", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

}

func TestRemoveItemFromCart_AbsentLine(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "DELETE", "/cart/items/"+testMenuItemID.String(), nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().RemoveItemFromCart(mock.Anything, testCustomerID, testMenuItemID).Return(domain.ErrCartNotFound)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestRemoveItemFromCart_Success(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "DELETE", "/cart/items/"+testMenuItemID.String(), nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().RemoveItemFromCart(mock.Anything, testCustomerID, testMenuItemID).Return(nil)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
