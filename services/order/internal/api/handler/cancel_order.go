package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.cancelOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := domain.OrderIDRequiredErrorCode
		log.DebugContext(r.Context(), msg, "error", orderID)
		h.respondJSON(w, http.StatusBadRequest, response.NewBaseError(msg))
		return
	}
	log.DebugContext(r.Context(), "order id parsed", "order_id", orderID)

	// orders, err := h.customer.CancelOrder(r.Context(), orderID)
	// if err != nil {
	// 	msg := "failed to cancel order"
	// 	log.Error(msg, "error", err)
	// 	h.respond(w, http.StatusInternalServerError, response.NewBaseError(msg, err))
	// 	return
	// }
	//

	h.respondJSON(w, http.StatusOK, response.Order{
		Items: []response.OrderItem{},
	})
}
