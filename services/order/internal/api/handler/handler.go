package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Storage interface {
}

// Handler holds all dependencies for HTTP handlers
type Handler struct {
	storage   Storage
	validator *validator.Validate
}

func New(storage Storage) http.Handler {
	h := &Handler{
		storage:   storage,
		validator: validator.New(),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/livez", h.livez)

	return mux
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}

func (h *Handler) livez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
