package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.removeItemFromCart"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	menuItemID, err := uuid.Parse(r.PathValue(menuItemIDKey))
	if err != nil {
		msg := common.MenuItemRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log = log.With("order_id", menuItemID)

	log.InfoContext(r.Context(), "request received")

	err = h.customer.RemoveItemFromCart(r.Context(), customerID, menuItemID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCartNotFound):
		default:
			log.ErrorContext(r.Context(), errInternalMsg, "error", err)
			common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError(errInternalMsg))
			return
		}
	}

	common.RespondJSON(log, w, http.StatusNoContent, nil)
}
