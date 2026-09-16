package handler

import (
	"errors"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.getCart"
	customerID := middleware.GetCustomerID(r.Context())
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), customerID,
	)

	log.InfoContext(r.Context(), "getting cart")

	cart, err := h.customer.GetCart(
		r.Context(),
		customerID,
	)
	if err != nil {
		if dErr, ok := errors.AsType[*domain.Error](err); ok {
			if dErr.Code == domain.CustomerNotFoundErrorCode {
				h.respondJSON(w, http.StatusNotFound, response.NewError(dErr))
				return
			}
		}

		log.ErrorContext(r.Context(), "failed to get cart", "error", err)
		h.respondJSON(w, http.StatusInternalServerError, nil)
		return
	}

	cartItems := make([]response.CardItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		cartItems = append(cartItems, response.CardItem{
			MenuItemID:     item.MenuItemID,
			Name:           item.Name,
			UnitPriceMinor: item.UnitPriceMinor,
			Currency:       item.Currency,
			Quantity:       item.Quantity,
			Instructions:   item.Instructions,
		})
	}

	h.respondJSON(w, http.StatusOK, response.Cart{
		RestaurantID:  cart.RestaurantID,
		Items:         cartItems,
		SubtotalMinor: cart.SubtotalMinor,
		Currency:      cart.Currency,
	})
}
