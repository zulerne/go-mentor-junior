package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
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
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(log, w, http.StatusBadRequest, response.NewBaseError(msg))
		return
	}

	log.InfoContext(r.Context(), "request received", "req", req)

	if err := h.validator.Struct(req); err != nil {
		msg := response.ValidationErrorCode
		log.ErrorContext(r.Context(), msg, "error", err)

		if validationErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			common.RespondJSON(log, w, http.StatusBadRequest, response.NewValidationError(validationErr))
			return
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		common.RespondJSON(log, w, http.StatusInternalServerError, response.NewBaseError("server error"))
		return
	}
	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	common.RespondJSON(log, w, http.StatusCreated, response.Order{
		Status:    string(domain.Pending),
		Items:     []response.OrderItem{},
		CreatedAt: date,
		UpdatedAt: date,
	})
}
