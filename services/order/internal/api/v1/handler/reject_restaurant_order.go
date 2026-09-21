package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type rejectRestaurantOrderRequest struct {
	Reason string `json:"reason" validate:"required"`
}

func (h *Handler) rejectRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.rejectRestaurantOrder"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	restaurantID, _ := middleware.RestaurantIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"restaurant_id", restaurantID,
	)

	var req rejectRestaurantOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		msg := "failed to decode request"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusBadRequest,
			common.NewBaseError(msg),
		)
		return
	}

	orderID, err := uuid.Parse(r.PathValue(orderIDKey))
	if err != nil {
		msg := common.OrderIDRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log = log.With("order_id", orderID)

	log.DebugContext(r.Context(), "request received", "reason", req.Reason)

	if err := h.validator.Struct(req); err != nil {
		msg := common.ValidationErrorCode
		log.ErrorContext(r.Context(), msg, "error", err)

		if _, ok := errors.AsType[validator.ValidationErrors](err); ok {
			common.RespondJSON(
				log,
				w,
				http.StatusBadRequest,
				common.NewError(common.RejectionReasonRequiredErrorCode, "reason is required", nil),
			)
			return
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusInternalServerError,
			common.NewBaseError("server error"),
		)
		return
	}

	order, err := h.restaurant.RejectOrder(r.Context(), restaurantID, orderID, req.Reason)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			common.RespondJSON(
				log,
				w,
				http.StatusNotFound,
				common.NewError(common.OrderNotFoundErrorCode, "order not found", nil),
			)
		case errors.Is(err, domain.ErrRestaurantIDMismatch):
			common.RespondJSON(
				log,
				w,
				http.StatusForbidden,
				common.NewError(common.OrderAccessDeniedErrorCode, "order access denied", nil),
			)
		case errors.Is(err, domain.ErrInvalidOrderStatus):
			common.RespondJSON(
				log,
				w,
				http.StatusConflict,
				common.NewError(common.InvalidOrderTransitionErrorCode, "invalid order status", nil),
			)
		default:
			msg := "failed to get order"
			log.ErrorContext(r.Context(), msg, "error", err)
			common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError(msg))
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
