package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
)

type Storage interface {
}

// Handler holds all dependencies for HTTP handlers
type Handler struct {
	storage   Storage
	validator *validator.Validate
	log       *slog.Logger
}

func New(storage Storage, log *slog.Logger) http.Handler {
	log = log.With("component", "handler")
	h := &Handler{
		storage:   storage,
		validator: validator.New(),
		log:       log,
	}

	customerIDMiddleware := middleware.CustomerID(log)
	restaurantIDMiddleware := middleware.RestaurantID(log)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.healthz)
	mux.HandleFunc("GET /livez", h.livez)

	mux.Handle("GET /orders", customerIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET /orders"))
	})))

	mux.Handle("GET /restaurants", restaurantIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET /restaurants"))
	})))

	return middleware.Chain(
		mux,
		middleware.RequestID(log),
		middleware.Logger(log),
	)
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (h *Handler) livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
