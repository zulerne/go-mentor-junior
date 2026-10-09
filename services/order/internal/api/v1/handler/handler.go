package handler

//go:generate go run github.com/vektra/mockery/v3@v3.8.0

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"time"

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
		item domain.CartItem,
	) (domain.Cart, error)
	RemoveItemFromCart(ctx context.Context, customerID uuid.UUID, menuItemID uuid.UUID) error
	CreateOrder(ctx context.Context, customerID uuid.UUID, deliveryAddress string) (domain.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (domain.Order, error)
	GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]domain.Order, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
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

type Database interface {
	Ping(ctx context.Context) error
}

// Handler holds all dependencies for HTTP handlers.
type Handler struct {
	pool       Database
	customer   Customer
	restaurant Restaurant
	validator  *validator.Validate
	log        *slog.Logger
}

func New(pool Database, cust Customer, rest Restaurant, validator *validator.Validate, log *slog.Logger) *Handler {
	log = log.With("component", "handler")
	h := &Handler{
		pool:       pool,
		customer:   cust,
		restaurant: rest,
		validator:  validator,
		log:        log,
	}

	return h
}

func (h *Handler) Routes() http.Handler {
	log := h.log
	mux := http.NewServeMux()

	// pprof endpoints
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)

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

func (h *Handler) readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	status := http.StatusOK

	if err := h.pool.Ping(ctx); err != nil {
		h.log.ErrorContext(ctx, "database ping failed", "error", err)
		status = http.StatusServiceUnavailable
	}

	common.RespondJSON(
		h.log,
		w,
		status,
		nil,
	)
}
