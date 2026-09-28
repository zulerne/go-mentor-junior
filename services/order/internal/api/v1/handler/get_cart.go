package handler

import (
	"errors"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.getCart"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	log.InfoContext(r.Context(), "getting cart")

	cart, err := h.customer.GetCart(
		r.Context(),
		customerID,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotFound):
			common.RespondJSON(
				log,
				w,
				http.StatusOK,
				CartResponse{
					Items: []CartItem{},
				},
			)
		default:
			msg := "failed to get cart"
			log.ErrorContext(r.Context(), msg, "error", err)
			common.RespondJSON(
				log,
				w,
				http.StatusInternalServerError,
				common.NewBaseError(msg),
			)
		}

		return
	}
	common.RespondJSON(log, w, http.StatusOK, cartToDTO(cart))
}
