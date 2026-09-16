package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type rejectRestaurantOrderRequest struct {
	Reason string `json:"reason"`
}

func (h *Handler) rejectRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.rejectRestaurantOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.RestaurantIDKey), middleware.GetRestaurantID(r.Context()),
	)

	var req rejectRestaurantOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request"
		log.DebugContext(r.Context(), msg, "error", err)
		h.respondJSON(w, http.StatusBadRequest, response.NewBaseError(msg))
		return
	}

	log.DebugContext(r.Context(), "request received", "reason", req.Reason)

	if err := h.validator.Struct(req); err != nil {
		msg := domain.ValidationErrorCode
		log.DebugContext(r.Context(), msg, "error", err)

		var validationErr validator.ValidationErrors
		if !errors.As(err, &validationErr) {
			h.respondJSON(w, http.StatusInternalServerError, response.NewBaseError(msg))
			return
		}

		h.respondJSON(w, http.StatusBadRequest, response.NewValidationError(validationErr))
		return
	}

	h.respondJSON(w, http.StatusOK, response.Order{
		Items: []response.OrderItem{},
	})
}
