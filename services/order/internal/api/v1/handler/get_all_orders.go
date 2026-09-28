package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
)

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	op := "handler.getAllOrders"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	orders, err := h.customer.GetAllOrders(r.Context(), customerID)
	if err != nil {
		msg := "failed to get all orders"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError(msg))
		return
	}

	common.RespondJSON(log, w, http.StatusOK, AllOrdersResponse{
		Orders: ordersToDTOs(orders),
	})
}
