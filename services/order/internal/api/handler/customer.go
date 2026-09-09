package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
)

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.getCart"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	log.Info("getting cart")

	// cart, err := h.orderService.GetCart(r.Context())
	// if err != nil {
	// 	log.Error("failed to get cart", "error", err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	h.respond(w, http.StatusOK, response.Cart{
		Items: []response.CardItem{},
	})
}

type addItemToCartRequest struct {
	RestaurantID string `json:"restaurant_id" validate:"required"`
	ItemID       string `json:"item_id" validate:"required"`
	Instructions string `json:"instructions" validate:"required"`
}

func (h *Handler) addItemToCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.addItemToCart"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	var req addItemToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request"
		log.Error(msg, "error", err)
		h.baseError(w, http.StatusBadRequest, msg)
		return
	}

	log.Info("request received")

	if err := h.validator.Struct(req); err != nil {
		msg := "validation error"
		log.Error(msg, "error", err)

		var validationErr validator.ValidationErrors
		if !errors.As(err, &validationErr) {
			h.baseError(w, http.StatusInternalServerError, "internal error")
			return
		}

		h.respond(w, http.StatusBadRequest, response.ValidationError(validationErr))
		return
	}

	// cart, err := h.cartService.AddItemToCart(r.Context(), req.ItemID, req.Quantity)
	// if err != nil {
	// 	msg := "failed to add item to cart"
	// 	log.Error(msg, "error", err)
	// 	h.baseError(w, http.StatusInternalServerError, msg)
	// 	return
	// }

	h.respond(w, http.StatusOK, response.Cart{
		Items: []response.CardItem{},
	})
}

func (h *Handler) removeItemFromCart(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {

}
