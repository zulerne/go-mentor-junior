package handler

import (
	"net/http"
	"time"

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

	responseOrders := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderItems := make([]OrderItem, 0, len(order.Items))
		for _, item := range order.Items {
			orderItems = append(orderItems, OrderItem{
				MenuItemID:     item.MenuItemID.String(),
				Name:           item.Name,
				UnitPriceMinor: item.UnitPriceMinor,
				Quantity:       item.Quantity,
				Instructions:   item.Instructions,
			})
		}
		responseOrders = append(responseOrders, OrderResponse{
			ID:              order.ID.String(),
			CustomerID:      order.CustomerID.String(),
			RestaurantID:    order.RestaurantID.String(),
			Status:          string(order.Status),
			Items:           orderItems,
			SubtotalMinor:   order.SubtotalMinor,
			Currency:        order.Currency,
			DeliveryAddress: order.DeliveryAddress,
			RejectionReason: order.RejectionReason,
			DeliveryStatus:  string(order.DeliveryStatus),
			CreatedAt:       order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
		})
	}

	common.RespondJSON(log, w, http.StatusOK, AllOrdersResponse{
		Orders: responseOrders,
	})
}
