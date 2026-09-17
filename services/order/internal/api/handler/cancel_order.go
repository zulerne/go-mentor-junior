package handler

import (
	"net/http"
	"time"

	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.cancelOrder"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"customer_id", middleware.GetCustomerID(r.Context()),
	)

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := response.OrderIDRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
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

	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	h.respondJSON(w, http.StatusOK, response.Order{
		Status:    string(domain.Cancelled),
		Items:     []response.OrderItem{},
		CreatedAt: date,
		UpdatedAt: date,
	})
}
