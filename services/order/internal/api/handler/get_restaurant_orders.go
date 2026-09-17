package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
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
			common.RespondJSON(
				log,
				w,
				http.StatusNotFound,
				common.NewError(common.RestaurantNotFoundErrorCode, "restaurant not found", nil),
			)
			return
		}
		msg := "failed to get orders"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusInternalServerError,
			common.NewBaseError(msg),
		)
		return
	}

	orderResponses := make([]response.Order, 0, len(orders))
	for _, order := range orders {
		var deliveryStatus *string
		if order.DeliveryStatus != nil {
			deliveryStatus = (*string)(order.DeliveryStatus)
		}

		items := make([]response.OrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			items = append(items, response.OrderItem{
				MenuItemID:     item.MenuItemID,
				Name:           item.Name,
				UnitPriceMinor: item.UnitPriceMinor,
				Quantity:       item.Quantity,
				Instructions:   item.Instructions,
			})
		}

		orderResponses = append(orderResponses, response.Order{
			ID:              order.ID,
			CustomerID:      order.CustomerID,
			RestaurantID:    order.RestaurantID,
			Status:          string(order.Status),
			SubtotalMinor:   order.SubtotalMinor,
			Currency:        order.Currency,
			DeliveryAddress: order.DeliveryAddress,
			RejectionReason: order.RejectionReason,
			DeliveryStatus:  deliveryStatus,
			CreatedAt:       order.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:       order.UpdatedAt.UTC().Format(time.RFC3339),
			Items:           items,
		})
	}

	common.RespondJSON(
		log,
		w,
		http.StatusOK,
		response.AllOrders{
			Orders: orderResponses,
		})
}
