package handler

//go:generate go run github.com/vektra/mockery/v3@v3.8.0

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

const (
	menuItemIDKey  = "menu_item_id"
	orderIDKey     = "order_id"
	errInternalMsg = "internal server error"
)

// Customer describes customer-facing operations.
//
//mockery:generate: true
type Customer interface {
	GetCart(ctx context.Context, customerID uuid.UUID) (domain.Cart, error)
	AddItemToCart(
		ctx context.Context,
		customerID uuid.UUID,
		restaurantID uuid.UUID,
		menuItemID uuid.UUID,
		quantity int32,
		instructions string,
	) (domain.Cart, error)
	RemoveItemFromCart(ctx context.Context, customerID uuid.UUID, menuItemID uuid.UUID) error
	CreateOrder(ctx context.Context, customerID uuid.UUID, deliveryAddress string) (domain.Order, error)
	GetOrder(ctx context.Context, customerID uuid.UUID, orderID uuid.UUID) (domain.Order, error)
	GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error)
	CancelOrder(ctx context.Context, customerID uuid.UUID, orderID uuid.UUID) (domain.Order, error)
}

// Restaurant describes restaurant-facing operations.
//
//mockery:generate: true
type Restaurant interface {
	GetOrders(ctx context.Context, restaurantID uuid.UUID) ([]domain.Order, error)
	AcceptOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (domain.Order, error)
	RejectOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID, reason string) (domain.Order, error)
	PrepareOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (domain.Order, error)
	ReadyOrder(ctx context.Context, restaurantID uuid.UUID, orderID uuid.UUID) (domain.Order, error)
}

// Handler holds all dependencies for HTTP handlers.
type Handler struct {
	customer   Customer
	restaurant Restaurant
	validator  *validator.Validate
	log        *slog.Logger
}

func New(cust Customer, rest Restaurant, validator *validator.Validate, log *slog.Logger) *Handler {
	log = log.With("component", "handler")
	h := &Handler{
		customer:   cust,
		restaurant: rest,
		validator:  validator,
		log:        log,
	}

	return h
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	log := h.log
	mux.HandleFunc("GET /livez", h.livez)
	mux.HandleFunc("GET /readyz", h.readyz)

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
	mux.Handle(
		"POST /restaurant/orders/{order_id}/accept",
		restaurantIDMiddleware(http.HandlerFunc(h.acceptRestaurantOrder)),
	)
	mux.Handle(
		"POST /restaurant/orders/{order_id}/reject",
		restaurantIDMiddleware(http.HandlerFunc(h.rejectRestaurantOrder)),
	)
	mux.Handle(
		"POST /restaurant/orders/{order_id}/start-preparation",
		restaurantIDMiddleware(http.HandlerFunc(h.prepareRestaurantOrder)),
	)
	mux.Handle(
		"POST /restaurant/orders/{order_id}/ready",
		restaurantIDMiddleware(http.HandlerFunc(h.readyRestaurantOrder)),
	)

	return middleware.Chain(
		mux,
		middleware.RequestID(log),
		middleware.Logger(log),
		middleware.Recoverer(log),
	)
}

func (h *Handler) livez(w http.ResponseWriter, _ *http.Request) {
	common.RespondJSON(
		h.log,
		w,
		http.StatusOK,
		nil,
	)
}

func (h *Handler) readyz(w http.ResponseWriter, _ *http.Request) {
	common.RespondJSON(
		h.log,
		w,
		http.StatusOK,
		nil,
	)
}
