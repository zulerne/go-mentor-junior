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

type rejectRestaurantOrderRequest struct {
	Reason string `json:"reason" validate:"required"`
}

func (h *Handler) rejectRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.rejectRestaurantOrder"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	restaurantID, _ := middleware.RestaurantIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"restaurant_id", restaurantID,
	)

	var req rejectRestaurantOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		msg := "failed to decode request"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusBadRequest,
			common.NewBaseError(msg),
		)
		return
	}

	log.DebugContext(r.Context(), "request received", "reason", req.Reason)

	if err := h.validator.Struct(req); err != nil {
		msg := common.ValidationErrorCode
		log.ErrorContext(r.Context(), msg, "error", err)

		if validationErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			common.RespondJSON(
				log,
				w,
				http.StatusBadRequest,
				common.NewValidationError(validationErr),
			)
			return
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusInternalServerError,
			common.NewBaseError("server error"),
		)
		return
	}

	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	common.RespondJSON(
		log,
		w,
		http.StatusOK,
		OrderResponse{
			Status:    string(domain.Rejected),
			Items:     []OrderItem{},
			CreatedAt: date,
			UpdatedAt: date,
		},
	)
}
