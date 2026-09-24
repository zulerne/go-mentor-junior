package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) prepareRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.prepareRestaurantOrder"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	restaurantID, _ := middleware.RestaurantIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"restaurant_id", restaurantID,
	)

	orderID, err := uuid.Parse(r.PathValue(orderIDKey))
	if err != nil {
		msg := common.OrderIDRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log = log.With("order_id", orderID)

	order, err := h.restaurant.PrepareOrder(r.Context(), restaurantID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			common.RespondJSON(
				log,
				w,
				http.StatusNotFound,
				common.NewError(common.OrderNotFoundErrorCode, "order not found", nil),
			)
		case errors.Is(err, domain.ErrOrderAccessDenied):
			common.RespondJSON(
				log,
				w,
				http.StatusForbidden,
				common.NewError(common.OrderAccessDeniedErrorCode, "order access denied", nil),
			)
		case errors.Is(err, domain.ErrInvalidOrderTransition):
			common.RespondJSON(
				log,
				w,
				http.StatusConflict,
				common.NewError(common.InvalidOrderTransitionErrorCode, "invalid order status", nil),
			)
		default:
			log.ErrorContext(r.Context(), errInternalMsg, "error", err)
			common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError(errInternalMsg))
		}

		return
	}

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

	common.RespondJSON(log, w, http.StatusOK, OrderResponse{
		ID:              orderID.String(),
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
