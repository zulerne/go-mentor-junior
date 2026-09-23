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

func TestAcceptRestaurantOrder_MissingRestaurantID(t *testing.T) {
	t.Parallel()

}

func TestAcceptRestaurantOrder_InvalidOrderID(t *testing.T) {
	t.Parallel()
}

func TestAcceptRestaurantOrder_NotFound(t *testing.T) {
	t.Parallel()
}

func TestAcceptRestaurantOrder_AccessDenied(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/accept", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().AcceptOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{}, domain.ErrOrderAccessDenied)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAcceptRestaurantOrder_InvalidTransition(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/accept", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().AcceptOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{}, domain.ErrInvalidOrderTransition)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestAcceptRestaurantOrder_DeliveryFailed(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/accept", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().AcceptOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{}, domain.ErrDeliveryProvider)

	hand.Routes().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestAcceptRestaurantOrder_Success(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	cust := NewMockCustomer(t)
	rest := NewMockRestaurant(t)
	hand := newHandler(cust, rest)

	req := httptest.NewRequestWithContext(context.Background(), "POST", "/restaurant/orders/"+testOrderID.String()+"/accept", nil)
	req.Header.Set(restaurantIDHeader, testRestaurantID.String())

	rest.EXPECT().AcceptOrder(mock.Anything, testRestaurantID, testOrderID).Return(domain.Order{
		ID:           testOrderID,
		CustomerID:   testCustomerID,
		RestaurantID: testRestaurantID,
		Status:       domain.Accepted,
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
	assert.Equal(t, string(domain.Accepted), resp.Status)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, testMenuItemID.String(), resp.Items[0].MenuItemID)
	assert.Equal(t, "Test Item", resp.Items[0].Name)
	assert.EqualValues(t, 1, resp.Items[0].Quantity)
	assert.EqualValues(t, 100, resp.Items[0].UnitPriceMinor)
	assert.Equal(t, "", resp.Items[0].Instructions)
}
