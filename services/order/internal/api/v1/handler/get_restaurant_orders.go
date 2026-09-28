package handler

import (
	"errors"
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

func (h *Handler) getRestaurantOrders(w http.ResponseWriter, r *http.Request) {
	op := "handler.getRestaurantOrders"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	restaurantID, _ := middleware.RestaurantIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"restaurant_id", restaurantID,
	)

	log.DebugContext(r.Context(), "request received")

	orders, err := h.restaurant.GetOrders(r.Context(), restaurantID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			common.RespondJSON(
				log,
				w,
				http.StatusNotFound,
				common.NewError(common.RestaurantNotFoundErrorCode, "restaurant not found", nil),
			)
			return
		}
		msg := "failed to get orders"
		log.ErrorContext(r.Context(), msg, "error", err)
		common.RespondJSON(
			log,
			w,
			http.StatusInternalServerError,
			common.NewBaseError(msg),
		)
		return
	}

	common.RespondJSON(log, w, http.StatusOK, AllOrdersResponse{
		Orders: ordersToDTOs(orders),
	})
}
