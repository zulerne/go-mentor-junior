package domain

type ErrorCode = string

const (
	BaseErrorCode            ErrorCode = "BASE_ERROR"
	ValidationErrorCode      ErrorCode = "VALIDATION_ERROR"
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
	Code    ErrorCode
	Message string
	Details []string
}

func (e Error) Error() string {
	return e.Message
}
