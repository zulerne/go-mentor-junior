package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
)

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	op := "handler.getAllOrders"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"customer_id", middleware.GetCustomerID(r.Context()),
	)

	// orders, err := h.customer.GetAllOrders(r.Context())
	// if err != nil {
	// 	msg := "failed to get all orders"
	// 	log.Error(msg, "error", err)
	// 	h.respond(w, http.StatusInternalServerError, response.NewBaseError(msg, err))
	// 	return
	// }
	//
	log.DebugContext(r.Context(), "getting all orders")

	common.RespondJSON(log, w, http.StatusOK, AllOrdersResponse{
		Orders: []OrderResponse{},
	})
}
