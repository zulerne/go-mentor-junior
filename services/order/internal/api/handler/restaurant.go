package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
)

func (h *Handler) getRestaurantOrders(w http.ResponseWriter, r *http.Request) {
	op := "handler.getRestaurantOrders"
	restaurantID := middleware.GetRestaurantID(r.Context())
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.RestaurantIDKey), restaurantID,
	)

	log.Debug("request received")

	orders, err := h.restaurant.GetOrders(r.Context(), restaurantID)
	if err != nil {
		h.respond(w, http.StatusInternalServerError, nil)
		return
	}

	orderResponses := make([]response.Order, 0, len(orders))
	for _, order := range orders {
		deliveryStatus := ""
		if order.DeliveryStatus != nil {
			deliveryStatus = string(*order.DeliveryStatus)
		}
		orderResponses = append(orderResponses, response.Order{
			ID:              order.ID,
			CustomerID:      order.CustomerID,
			Status:          string(order.Status),
			SubtotalMinor:   order.SubtotalMinor,
			Currency:        order.Currency,
			DeliveryAddress: order.DeliveryAddress,
			RejectionReason: order.RejectionReason,
			DeliveryStatus:  &deliveryStatus,
			CreatedAt:       order.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       order.UpdatedAt.Format(time.RFC3339),
			Items:           []response.OrderItem{},
		})
	}

	h.respond(w, http.StatusOK, response.AllOrders{
		Orders: orderResponses,
	})
}

func (h *Handler) acceptRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.acceptRestaurantOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.RestaurantIDKey), middleware.GetRestaurantID(r.Context()),
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

type rejectRestaurantOrderRequest struct {
	Reason string `json:"reason"`
}

func (h *Handler) rejectRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.rejectRestaurantOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.RestaurantIDKey), middleware.GetRestaurantID(r.Context()),
	)

	var req rejectRestaurantOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := "failed to decode request"
		log.Error(msg, "error", err)
		h.respond(w, http.StatusBadRequest, response.NewBaseError(msg, nil))
		return
	}

	log.Debug("request received", "reason", req.Reason)

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

	h.respond(w, http.StatusOK, response.Order{
		Items: []response.OrderItem{},
	})
}

func (h *Handler) prepareRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.prepareRestaurantOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.RestaurantIDKey), middleware.GetRestaurantID(r.Context()),
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

func (h *Handler) readyRestaurantOrder(w http.ResponseWriter, r *http.Request) {
	op := "handler.readyRestaurantOrder"
	log := h.log.With(
		"op", op,
		string(middleware.RequestIDKey), middleware.GetRequestID(r.Context()),
		string(middleware.RestaurantIDKey), middleware.GetRestaurantID(r.Context()),
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
