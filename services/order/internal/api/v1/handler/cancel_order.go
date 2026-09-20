package handler

import (
	"net/http"
	"time"

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

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := common.OrderIDRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
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
	common.RespondJSON(log, w, http.StatusOK, OrderResponse{
		Status:    string(domain.Cancelled),
		Items:     []OrderItem{},
		CreatedAt: date,
		UpdatedAt: date,
	})
}
