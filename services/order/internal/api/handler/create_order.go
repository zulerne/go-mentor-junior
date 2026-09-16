package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
)

type createOrderRequest struct {
	DeliveryAddress string `json:"delivery_address" validate:"required"`
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.createOrder"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"customer_id", middleware.GetCustomerID(r.Context()),
	)

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request body"
		log.DebugContext(r.Context(), msg, "error", err)
		h.respondJSON(w, http.StatusBadRequest, response.NewBaseError(msg))
		return
	}

	log.InfoContext(r.Context(), "request received", "req", req)

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

	h.respondJSON(w, http.StatusCreated, response.Order{
		Items: []response.OrderItem{},
	})
}
