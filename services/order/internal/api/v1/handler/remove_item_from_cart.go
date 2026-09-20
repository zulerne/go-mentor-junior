package handler

import (
	"net/http"

	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
)

func (h *Handler) removeItemFromCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.removeItemFromCart"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	menuItemID := r.PathValue(menuItemIDKey)
	if menuItemID == "" {
		msg := "menu item id is required"
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(
			log,
			w,
			http.StatusBadRequest,
			common.NewBaseError(msg),
		)
		return
	}
	log.DebugContext(r.Context(), "menu item id parsed", "menu_item_id", menuItemID)

	common.RespondJSON(
		log,
		w,
		http.StatusNoContent,
		nil,
	)
}
