package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.getOrder"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"customer_id", middleware.GetCustomerID(r.Context()),
	)

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := domain.OrderIDRequiredErrorCode
		log.DebugContext(r.Context(), msg, "error", orderID)
		h.respondJSON(w, http.StatusBadRequest, response.NewBaseError(msg))
		return
	}
	log.DebugContext(r.Context(), "order id parsed", "order_id", orderID)

	h.respondJSON(w, http.StatusOK, response.Order{
		Items: []response.OrderItem{},
	})
}
