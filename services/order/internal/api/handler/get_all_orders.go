package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
)

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
	log.DebugContext(r.Context(), "getting all orders")

	h.respondJSON(w, http.StatusOK, response.AllOrders{
		Orders: []response.Order{},
	})
}
