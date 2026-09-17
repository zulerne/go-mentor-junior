package handler

import (
	"net/http"
	"time"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) acceptRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.acceptRestaurantOrder"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"restaurant_id", middleware.GetRestaurantID(r.Context()),
	)

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := common.OrderIDRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log.DebugContext(r.Context(), "order id parsed", "order_id", orderID)

	date := time.Date(2026, time.September, 1, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	common.RespondJSON(log, w, http.StatusOK, response.Order{
		Status:    string(domain.Accepted),
		Items:     []response.OrderItem{},
		CreatedAt: date,
		UpdatedAt: date,
	})
}
