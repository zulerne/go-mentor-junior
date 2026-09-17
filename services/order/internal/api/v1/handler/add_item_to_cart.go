package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
)

type addItemToCartRequest struct {
	RestaurantID string `json:"restaurant_id" validate:"required"`
	Quantity     int    `json:"quantity"      validate:"required"`
	Instructions string `json:"instructions"`
}

func (h *Handler) addItemToCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.addItemToCart"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"customer_id", middleware.GetCustomerID(r.Context()),
	)

	menuItemID := r.PathValue(menuItemIDKey)
	if menuItemID == "" {
		msg := "menu item id is required"
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log.DebugContext(r.Context(), "menu item id parsed", "menu_item_id", menuItemID)

	var req addItemToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}

	log.InfoContext(r.Context(), "request received")

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

	// cart, err := h.cartService.AddItemToCart(r.Context(), req.ItemID, req.Quantity)
	// if err != nil {
	// 	msg := "failed to add item to cart"
	// 	log.Error(msg, "error", err)
	// 	h.baseError(w, http.StatusInternalServerError, msg)
	// 	return
	// }

	common.RespondJSON(log, w, http.StatusOK, CartResponse{
		Items: []CartItem{},
	})
}
