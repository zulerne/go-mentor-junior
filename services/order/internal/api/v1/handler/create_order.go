package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type createOrderRequest struct {
	DeliveryAddress string `json:"delivery_address" validate:"required,max=500"`
}

//nolint:funlen
func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.createOrder"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	var req createOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		msg := "failed to decode request body"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}

	log.InfoContext(r.Context(), "request received", "req", req)

	req.DeliveryAddress = strings.TrimSpace(req.DeliveryAddress)

	if err := h.validator.Struct(req); err != nil {
		msg := common.ValidationErrorCode
		log.ErrorContext(r.Context(), msg, "error", err)

		if validationErr, ok := errors.AsType[validator.ValidationErrors](err); ok {
			common.RespondJSON(log, w, http.StatusBadRequest, common.NewValidationError(validationErr))
			return
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError("server error"))
		return
	}

	if req.DeliveryAddress == "" {
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewError(common.InvalidDeliveryAddressErrorCode, "delivery address is required", nil))
		return
	}

	order, err := h.customer.CreateOrder(r.Context(), customerID, req.DeliveryAddress)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCartEmpty):
			common.RespondJSON(log, w, http.StatusUnprocessableEntity, common.NewError(
				common.EmptyCartErrorCode,
				"cart empty",
				nil))
		case errors.Is(err, domain.ErrMenuItemNotAvailable):
			common.RespondJSON(log, w, http.StatusUnprocessableEntity, common.NewError(
				common.MenuItemUnavailableErrorCode,
				"menu item not available",
				nil))
		case errors.Is(err, domain.ErrInvalidDeliveryAddress):
			common.RespondJSON(log, w, http.StatusBadRequest, common.NewError(
				common.InvalidDeliveryAddressErrorCode,
				"invalid delivery address",
				nil))
		case errors.Is(err, domain.ErrRestaurantNotFound):
			common.RespondJSON(log, w, http.StatusNotFound, common.NewError(
				common.RestaurantNotFoundErrorCode,
				"restaurant not found",
				nil))
		case errors.Is(err, domain.ErrRestaurantNotAcceptingOrders):
			common.RespondJSON(log, w, http.StatusBadRequest, common.NewError(
				common.RestaurantNotAcceptingOrdersErrorCode,
				"restaurant not accepting orders",
				nil))
		case errors.Is(err, domain.ErrMinOrderNotReached):
			common.RespondJSON(log, w, http.StatusBadRequest, common.NewError(
				common.MinimumOrderNotReachedErrorCode,
				"min order amount not met",
				nil))
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

	common.RespondJSON(log, w, http.StatusCreated, OrderResponse{
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
