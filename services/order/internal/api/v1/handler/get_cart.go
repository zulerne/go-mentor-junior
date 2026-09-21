package handler

import (
	"errors"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.getCart"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	log.InfoContext(r.Context(), "getting cart")

	cart, err := h.customer.GetCart(
		r.Context(),
		customerID,
	)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			common.RespondJSON(
				log,
				w,
				http.StatusNotFound,
				common.NewError(common.CustomerNotFoundErrorCode, "customer not found", nil),
			)
			return
		}

		msg := "failed to get cart"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusInternalServerError,
			common.NewBaseError(msg),
		)
		return
	}

	cartItems := make([]CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		cartItems = append(cartItems, CartItem{
			MenuItemID:     item.MenuItemID.String(),
			Name:           item.Name,
			UnitPriceMinor: item.UnitPriceMinor,
			Currency:       item.Currency,
			Quantity:       item.Quantity,
			Instructions:   item.Instructions,
		})
	}

	common.RespondJSON(
		log,
		w,
		http.StatusOK,
		CartResponse{
			RestaurantID:  cart.RestaurantID.String(),
			Items:         cartItems,
			SubtotalMinor: cart.SubtotalMinor,
			Currency:      cart.Currency,
		},
	)
}
