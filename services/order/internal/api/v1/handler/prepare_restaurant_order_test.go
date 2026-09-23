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

func TestPrepareRestaurantOrder_MissingRestaurantID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/start-preparation", nil)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPrepareRestaurantOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+"invalid"+"/start-preparation", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestPrepareRestaurantOrder_NotFound(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/start-preparation", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().PrepareOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{}, domain.ErrOrderNotFound)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPrepareRestaurantOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/start-preparation", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().PrepareOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{}, domain.ErrOrderAccessDenied)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestPrepareRestaurantOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/start-preparation", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().PrepareOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{}, domain.ErrInvalidOrderTransition)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestPrepareRestaurantOrder_Success(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/start-preparation", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().PrepareOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Preparing,
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

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	resp := decodeOrder(t, rec)
	assert.Equal(t, testOrderID.String(), resp.ID)
	assert.Equal(t, testCustomerID.String(), resp.CustomerID)
	assert.Equal(t, testRestaurantID.String(), resp.RestaurantID)
	assert.Equal(t, string(domain.Preparing), resp.Status)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testMenuItemID.String(), resp.Items[0].MenuItemID)
	assert.Equal(t, "Test Item", resp.Items[0].Name)
	assert.EqualValues(t, 1, resp.Items[0].Quantity)
	assert.EqualValues(t, 100, resp.Items[0].UnitPriceMinor)
	assert.Equal(t, "", resp.Items[0].Instructions)
}
