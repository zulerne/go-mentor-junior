package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type rejectRestaurantOrderRequest struct {
	Reason string `json:"reason" validate:"required"`
}

func (h *Handler) rejectRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.rejectRestaurantOrder"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"restaurant_id", middleware.GetRestaurantID(r.Context()),
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
		msg := response.ValidationErrorCode
		log.DebugContext(r.Context(), msg, "error", err)

		if validationErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			h.respondJSON(w, http.StatusBadRequest, response.NewValidationError(validationErr))
			return
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		h.respondJSON(w, http.StatusInternalServerError, response.NewBaseError("server error"))
		return
	}

	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	h.respondJSON(w, http.StatusOK, response.Order{
		Status:    string(domain.Rejected),
		Items:     []response.OrderItem{},
		CreatedAt: date,
		UpdatedAt: date,
	})
}
