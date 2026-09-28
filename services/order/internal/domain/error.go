package domain

import "errors"

var ErrNotFound = errors.New("not found")

var ErrInvalidOrderTransition = errors.New("invalid order status")
var ErrOrderAccessDenied = errors.New("order access denied")

var ErrRestaurantNotAcceptingOrders = errors.New("restaurant not accepting orders")
var ErrMenuItemNotAvailable = errors.New("item not available")
var ErrMinOrderNotReached = errors.New("min order amount not met")

var ErrDeliveryProvider = errors.New("delivery provider error")
var ErrInvalidDeliveryAddress = errors.New("invalid delivery address")

var ErrCartLimit = errors.New("invalid quantity")
var ErrCartEmpty = errors.New("cart empty")
