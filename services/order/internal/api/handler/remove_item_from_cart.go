package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/api/response"
)

func (h *Handler) removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.removeItemFromCart"
	log := h.log.With(
		"op", op,
		"request_id", middleware.GetRequestID(r.Context()),
		"customer_id", middleware.GetCustomerID(r.Context()),
	)

	menuItemID := r.PathValue(menuItemIDKey)
	if menuItemID == "" {
		msg := "menu item id is required"
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(
			log,
			w,
			http.StatusBadRequest,
			response.NewBaseError(msg),
		)
		return
	}
	log.DebugContext(r.Context(), "menu item id parsed", "menu_item_id", menuItemID)

	common.RespondJSON(
		h.log,
		w,
		http.StatusNoContent,
		nil,
	)
}
