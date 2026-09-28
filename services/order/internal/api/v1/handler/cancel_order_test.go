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

func TestCancelOrder_ServiceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		serviceErr error
		wantStatus int
	}{
		{"not found", domain.ErrNotFound, http.StatusNotFound},
		{"invalid transition", domain.ErrInvalidOrderTransition, http.StatusConflict},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			cust := NewMockCustomer(t)
			handler := newHandler(cust, NewMockRestaurant(t))

			req := httptest.NewRequestWithContext(
				context.Background(),
				"POST",
				"/orders/"+testOrderID.String()+"/cancel",
				nil,
			)
			req.Header.Set(customerIDHeader, testCustomerID.String())

			cust.EXPECT().CancelOrder(mock.Anything, testOrderID).Return(tc.serviceErr)

			handler.Routes().ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}

func TestCancelOrder_Success(t *testing.T) {
	t.Parallel()
	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	handler := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/orders/"+testOrderID.String()+"/cancel", nil)
	req.Header.Set(customerIDHeader, testCustomerID.String())

	cust.EXPECT().CancelOrder(mock.Anything, testOrderID).Return(nil)

	handler.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
