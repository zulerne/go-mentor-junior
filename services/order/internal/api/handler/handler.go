package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/services/customer"
	"github.com/zulerne/go-mentor-junior/order/internal/services/restaurant"
)

// Handler holds all dependencies for HTTP handlers
type Handler struct {
	customer   *customer.Service
	restaurant *restaurant.Service
	validator  *validator.Validate
	log        *slog.Logger
}

func New(log *slog.Logger) http.Handler {
	log = log.With("component", "handler")
	h := &Handler{
		validator: validator.New(),
		log:       log,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.healthz)
	mux.HandleFunc("GET /livez", h.livez)

	customerIDMiddleware := middleware.CustomerID(log)

	mux.Handle("GET /cart", customerIDMiddleware(http.HandlerFunc(h.getCart)))
	mux.Handle("PUT /cart/items/{menu_item_id}", customerIDMiddleware(http.HandlerFunc(h.addItemToCart)))
	mux.Handle("DELETE /cart/items/{menu_item_id}", customerIDMiddleware(http.HandlerFunc(h.removeItemFromCart)))
	mux.Handle("POST /orders", customerIDMiddleware(http.HandlerFunc(h.createOrder)))
	mux.Handle("GET /orders/{order_id}", customerIDMiddleware(http.HandlerFunc(h.getOrder)))
	mux.Handle("GET /orders", customerIDMiddleware(http.HandlerFunc(h.getAllOrders)))
	mux.Handle("POST /orders/{order_id}/cancel", customerIDMiddleware(http.HandlerFunc(h.cancelOrder)))

	restaurantIDMiddleware := middleware.RestaurantID(log)

	mux.Handle("GET /restaurant/orders", restaurantIDMiddleware(http.HandlerFunc(h.getRestaurantOrders)))
	mux.Handle("POST /restaurant/orders/{order_id}/accept", restaurantIDMiddleware(http.HandlerFunc(h.acceptRestaurantOrder)))
	mux.Handle("POST /restaurant/orders/{order_id}/reject", restaurantIDMiddleware(http.HandlerFunc(h.rejectRestaurantOrder)))
	mux.Handle("POST /restaurant/orders/{order_id}/start-preparation", restaurantIDMiddleware(http.HandlerFunc(h.prepareRestaurantOrder)))
	mux.Handle("POST /restaurant/orders/{order_id}/ready", restaurantIDMiddleware(http.HandlerFunc(h.readyRestaurantOrder)))

	return middleware.Chain(
		mux,
		middleware.RequestID(log),
		middleware.Logger(log),
	)
}

func (h *Handler) respond(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) error(w http.ResponseWriter, status int, code string, message string, details any) {
	h.respond(w, status, response.Error{
		ErrorData: response.ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func (h *Handler) baseError(w http.ResponseWriter, status int, err string) {
	h.error(w, status, "INTERNAL_ERROR", err, nil)
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (h *Handler) livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
