package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/v1/handler"
)

const (
	customerIDHeader   = "X-Customer-ID"
	restaurantIDHeader = "X-Restaurant-ID"
)

var (
	testCustomerID   = uuid.MustParse("ef55f77a-7738-426f-93a4-f78f2baf7970")
	testRestaurantID = uuid.MustParse("fe70cb38-d10d-452c-8860-1af5936f7037")
	testOrderID      = uuid.MustParse("2f5efdc6-c936-4fdf-a681-1e21e73e6e71")
	testMenuItemID   = uuid.MustParse("58185771-1f9f-4d71-9609-bdacea1deb2e")
)

func newHandler(cust handler.Customer, rest handler.Restaurant) *handler.Handler {
	log := slog.New(slog.DiscardHandler)
	return handler.New(cust, rest, validator.New(), log)
}

func decodeCart(t *testing.T, rec *httptest.ResponseRecorder) handler.CartResponse {
	t.Helper()
	var resp handler.CartResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode CartResponse: %v", err)
	}
	return resp
}

func decodeOrder(t *testing.T, rec *httptest.ResponseRecorder) handler.OrderResponse {
	t.Helper()
	var resp handler.OrderResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode OrderResponse: %v", err)
	}
	return resp
}

func decodeAllOrders(t *testing.T, rec *httptest.ResponseRecorder) handler.AllOrdersResponse {
	t.Helper()
	var resp handler.AllOrdersResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode AllOrdersResponse: %v", err)
	}
	return resp
}

func jsonBody(t *testing.T, v any) io.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	return bytes.NewReader(b)
}
