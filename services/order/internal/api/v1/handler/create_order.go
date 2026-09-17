package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
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
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		msg := "failed to decode request body"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}

	log.InfoContext(r.Context(), "request received", "req", req)

	if err := h.validator.Struct(req); err != nil {
		msg := common.ValidationErrorCode
		log.ErrorContext(r.Context(), msg, "error", err)

		if validationErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			common.RespondJSON(log, w, http.StatusBadRequest, common.NewValidationError(validationErr))
			return
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError("server error"))
		return
	}
	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	common.RespondJSON(log, w, http.StatusCreated, OrderResponse{
		Status:    string(domain.Pending),
		Items:     []OrderItem{},
		CreatedAt: date,
		UpdatedAt: date,
	})
}
