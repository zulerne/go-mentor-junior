package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.cancelOrder"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	orderID, err := uuid.Parse(r.PathValue(orderIDKey))
	if err != nil {
		msg := common.OrderIDRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log = log.With("order_id", orderID)

	err = h.customer.CancelOrder(r.Context(), orderID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			common.RespondJSON(
				log,
				w,
				http.StatusNotFound,
				common.NewError(common.OrderNotFoundErrorCode, "order not found", nil),
			)
		case errors.Is(err, domain.ErrOrderAccessDenied):
			common.RespondJSON(
				log,
				w,
				http.StatusForbidden,
				common.NewError(common.OrderAccessDeniedErrorCode, "order access denied", nil),
			)
		case errors.Is(err, domain.ErrInvalidOrderTransition):
			common.RespondJSON(
				log,
				w,
				http.StatusConflict,
				common.NewError(common.InvalidOrderTransitionErrorCode, "invalid order transition", nil),
			)
		default:
			log.ErrorContext(r.Context(), errInternalMsg, "error", err)
			common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError(errInternalMsg))
		}

		return
	}

	common.RespondJSON(log, w, http.StatusNoContent, nil)
}
