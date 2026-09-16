package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) getRestaurantOrders(w http.ResponseWriter, r *http.Request) {
	op := "handler.getRestaurantOrders"
	restaurantID := middleware.GetRestaurantID(r.Context())
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"restaurant_id", restaurantID,
	)

	log.DebugContext(r.Context(), "request received")

	orders, err := h.restaurant.GetOrders(r.Context(), restaurantID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.respondJSON(
				w,
				http.StatusNotFound,
				response.NewError(response.RestaurantNotFoundErrorCode, "restaurant not found", nil),
			)
			return
		}
		msg := "failed to get orders"
		log.ErrorContext(r.Context(), msg, "error", err)
		h.respondJSON(w, http.StatusInternalServerError, response.NewBaseError(msg))
		return
	}

	orderResponses := make([]response.Order, 0, len(orders))
	for _, order := range orders {
		deliveryStatus := ""
		if order.DeliveryStatus != nil {
			deliveryStatus = string(*order.DeliveryStatus)
		}
		orderResponses = append(orderResponses, response.Order{
			ID:              order.ID,
			CustomerID:      order.CustomerID,
			Status:          string(order.Status),
			SubtotalMinor:   order.SubtotalMinor,
			Currency:        order.Currency,
			DeliveryAddress: order.DeliveryAddress,
			RejectionReason: order.RejectionReason,
			DeliveryStatus:  &deliveryStatus,
			CreatedAt:       order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
			Items:           []response.OrderItem{},
		})
	}

	h.respondJSON(w, http.StatusOK, response.AllOrders{
		Orders: orderResponses,
	})
}
