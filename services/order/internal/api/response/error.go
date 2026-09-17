package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorCode = string

const (
	BaseErrorCode            ErrorCode = "BASE_ERROR"
	ValidationErrorCode      ErrorCode = "VALIDATION_ERROR"
	BadRequestErrorCode      ErrorCode = "BAD_REQUEST"
	OrderIDRequiredErrorCode ErrorCode = "ORDER_ID_REQUIRED"

	CustomerNotFoundErrorCode             ErrorCode = "CUSTOMER_NOT_FOUND"
	RestaurantNotFoundErrorCode           ErrorCode = "RESTAURANT_NOT_FOUND"
	RestaurantNotAcceptingOrdersErrorCode ErrorCode = "RESTAURANT_NOT_ACCEPTING_ORDERS"
	MenuItemNotFoundErrorCode             ErrorCode = "MENU_ITEM_NOT_FOUND"
	MenuItemUnavailableErrorCode          ErrorCode = "MENU_ITEM_UNAVAILABLE"
	CartRestaurantConflictErrorCode       ErrorCode = "CART_RESTAURANT_CONFLICT"
	CartLimitExceededErrorCode            ErrorCode = "CART_LIMIT_EXCEEDED"
	InvalidQuantityErrorCode              ErrorCode = "INVALID_QUANTITY"
	InvalidInstructionsErrorCode          ErrorCode = "INVALID_INSTRUCTIONS"
	EmptyCartErrorCode                    ErrorCode = "EMPTY_CART"
	MinimumOrderNotReachedErrorCode       ErrorCode = "MINIMUM_ORDER_NOT_REACHED"
	InvalidDeliveryAddressErrorCode       ErrorCode = "INVALID_DELIVERY_ADDRESS"
	OrderNotFoundErrorCode                ErrorCode = "ORDER_NOT_FOUND"
	OrderAccessDeniedErrorCode            ErrorCode = "ORDER_ACCESS_DENIED"
	RejectionReasonRequiredErrorCode      ErrorCode = "REJECTION_REASON_REQUIRED"
	InvalidOrderTransitionErrorCode       ErrorCode = "INVALID_ORDER_TRANSITION"
	DeliveryCreationFailedErrorCode       ErrorCode = "DELIVERY_CREATION_FAILED"
	DeliveryTransitionFailedErrorCode     ErrorCode = "DELIVERY_TRANSITION_FAILED"
	DependencyUnavailableErrorCode        ErrorCode = "DEPENDENCY_UNAVAILABLE"
)

type Error struct {
	ErrorData ErrorData `json:"error"`
}

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func NewError(code ErrorCode, msg string, details any) Error {
	return Error{
		ErrorData: ErrorData{
			Code:    code,
			Message: msg,
			Details: details,
		},
	}
}

func NewBaseError(msg string) Error {
	return Error{
		ErrorData: ErrorData{
			Code:    BaseErrorCode,
			Message: msg,
			Details: map[string]any{},
		},
	}
}

// NewValidationError creates a new validation error from the given validator.ValidationErrors ().
func NewValidationError(errs validator.ValidationErrors) Error {
	var msgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("'%s' is required", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("'%s' is invalid", err.Field()))
		}
	}

	return Error{
		ErrorData: ErrorData{
			Code:    ValidationErrorCode,
			Message: strings.Join(msgs, ", "),
			Details: map[string]any{},
		},
	}
}
