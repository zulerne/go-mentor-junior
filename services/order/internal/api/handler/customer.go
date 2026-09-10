package handler

// TODO (review): is it okay that I splitted the code into customer/restaurant files?
import (
	"encoding/json"
	"errors"
	"fmt"
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
	fmt.Printf("h.Customer: %v\n", h.Customer)

	cart, err := h.Customer.CartStore.FindCart(
		r.Context(),
		middleware.GetCustomerID(r.Context()),
	)
	if err != nil {
		log.Error("failed to get cart", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respond(w, http.StatusOK, cart)
}

type addItemToCartRequest struct {
	RestaurantID string `json:"restaurant_id" validate:"required"`
	Quantity     int    `json:"quantity" validate:"required"`
	Instructions string `json:"instructions" validate:"required"`
}

func (h *Handler) addItemToCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.addItemToCart"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	menuItemID := r.PathValue(menuItemIDKey)
	if menuItemID == "" {
		msg := "menu item id is required"
		log.Error(msg, "error", menuItemID)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, nil))
		return
	}
	log.Debug("menu item id parsed", "menu_item_id", menuItemID)

	var req addItemToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request"
		log.Error(msg, "error", err)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, err))
		return
	}

	log.Info("request received")

	if err := h.validator.Struct(req); err != nil {
		msg := "validation error"
		log.Error(msg, "error", err)

		var validationErr validator.ValidationErrors
		if !errors.As(err, &validationErr) {
			h.respond(w, http.StatusInternalServerError, response.NewBaseError(msg, err))
			return
		}

		h.respond(w, http.StatusBadRequest, response.NewValidationError(validationErr))
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
	op := "handler.removeItemFromCart"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	menuItemID := r.PathValue(menuItemIDKey)
	if menuItemID == "" {
		msg := "menu item id is required"
		log.Error(msg, "error", menuItemID)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, nil))
		return
	}
	log.Debug("menu item id parsed", "menu_item_id", menuItemID)

	h.respond(w, http.StatusNoContent, nil)
}

type createOrderRequest struct {
	DeliveryAddress string `json:"delivery_address" validate:"required"`
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.createOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request body"
		log.Error(msg, "error", err)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, err))
		return
	}

	log.Info("request received", "req", req)

	if err := h.validator.Struct(req); err != nil {
		msg := "validation error"
		log.Error(msg, "error", err)

		var validationErr validator.ValidationErrors
		if !errors.As(err, &validationErr) {
			h.respond(w, http.StatusInternalServerError, response.NewBaseError(msg, err))
			return
		}

		h.respond(w, http.StatusBadRequest, response.NewValidationError(validationErr))
		return
	}

	h.respond(w, http.StatusCreated, response.Order{
		Items: []response.OrderItem{},
	})
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.getOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := "order id is required"
		log.Error(msg, "error", orderID)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, nil))
		return
	}
	log.Debug("order id parsed", "order_id", orderID)

	h.respond(w, http.StatusOK, response.Order{
		Items: []response.OrderItem{},
	})
}

func (h *Handler) getAllOrders(w http.ResponseWriter, r *http.Request) {
	op := "handler.getAllOrders"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	// orders, err := h.customer.GetAllOrders(r.Context())
	// if err != nil {
	// 	msg := "failed to get all orders"
	// 	log.Error(msg, "error", err)
	// 	h.respond(w, http.StatusInternalServerError, response.NewBaseError(msg, err))
	// 	return
	// }
	//
	log.Debug("getting all orders")

	h.respond(w, http.StatusOK, response.AllOrders{
		Orders: []response.Order{},
	})
}

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.cancelOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.CustomerIDKey), middleware.GetCustomerID(r.Context()),
	)

	orderID := r.PathValue(orderIDKey)
	if orderID == "" {
		msg := "order id is required"
		log.Error(msg, "error", orderID)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, nil))
		return
	}
	log.Debug("order id parsed", "order_id", orderID)

	// orders, err := h.customer.CancelOrder(r.Context(), orderID)
	// if err != nil {
	// 	msg := "failed to cancel order"
	// 	log.Error(msg, "error", err)
	// 	h.respond(w, http.StatusInternalServerError, response.NewBaseError(msg, err))
	// 	return
	// }
	//

	h.respond(w, http.StatusOK, response.Order{
		Items: []response.OrderItem{},
	})
}
