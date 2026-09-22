package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/api/common"
	"github.com/zulerne/go-mentor-junior/order/internal/api/middleware"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type addItemToCartRequest struct {
	RestaurantID uuid.UUID `json:"restaurant_id" validate:"required"`
	Quantity     int       `json:"quantity"      validate:"required,min=1,max=10"`
	Instructions string    `json:"instructions"  validate:"omitempty,max=250"`
}

func (h *Handler) addItemToCart(w http.ResponseWriter, r *http.Request) {
	op := "handler.addItemToCart"
	requestID, _ := middleware.RequestIDFromContext(r.Context())
	customerID, _ := middleware.CustomerIDFromContext(r.Context())

	log := h.log.With(
		"op", op,
		"request_id", requestID,
		"customer_id", customerID,
	)

	menuItemID, err := uuid.Parse(r.PathValue(menuItemIDKey))
	if err != nil {
		msg := common.MenuItemRequiredErrorCode
		log.ErrorContext(r.Context(), msg)
		common.RespondJSON(log, w, http.StatusBadRequest, common.NewBaseError(msg))
		return
	}
	log = log.With("order_id", menuItemID)

	var req addItemToCartRequest
	// TODO: common decoder?
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
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
			for _, fe := range validationErr {
				switch fe.Field() {
				case "Instructions":
					common.RespondJSON(log, w, http.StatusBadRequest, common.NewError(common.InvalidInstructionsErrorCode, "invalid instructions", nil))
					return
				case "Quantity":
					common.RespondJSON(log, w, http.StatusBadRequest, common.NewError(common.InvalidQuantityErrorCode, "invalid quantity", nil))
					return
				default:
					common.RespondJSON(log, w, http.StatusBadRequest, common.NewValidationError(validationErr))
					return
				}
			}
		}

		log.ErrorContext(r.Context(), "failed to validate request", "error", err)
		common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError("server error"))
		return
	}

	cart, err := h.customer.AddItemToCart(r.Context(), customerID, req.RestaurantID, menuItemID, req.Quantity, req.Instructions)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCartLimit):
			common.RespondJSON(
				log,
				w,
				http.StatusUnprocessableEntity,
				common.NewError(common.CartLimitExceededErrorCode, "invalid quantity", nil),
			)
		case errors.Is(err, domain.ErrRestaurantIDMismatch):
			common.RespondJSON(
				log,
				w,
				http.StatusConflict,
				common.NewError(common.CartRestaurantConflictErrorCode, "order access denied", nil),
			)
		case errors.Is(err, domain.ErrMenuItemNotAvailable):
			common.RespondJSON(
				log,
				w,
				http.StatusUnprocessableEntity,
				common.NewError(common.MenuItemUnavailableErrorCode, "item not available", nil),
			)
		default:
			msg := "failed to get order"
			log.ErrorContext(r.Context(), msg, "error", err)
			common.RespondJSON(log, w, http.StatusInternalServerError, common.NewBaseError(msg))
		}

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

	common.RespondJSON(log, w, http.StatusOK, CartResponse{
		RestaurantID:  cart.RestaurantID.String(),
		Items:         cartItems,
		SubtotalMinor: cart.SubtotalMinor,
		Currency:      cart.Currency,
	})
}
