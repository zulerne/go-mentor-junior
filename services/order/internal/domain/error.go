package domain

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrInvalidOrderStatus = errors.New("invalid order status")
var ErrOrderAccessDenied = errors.New("order access denied")
var ErrOrderAlreadyCancelled = errors.New("order already cancelled")

var ErrCustomerNotFound = errors.New("customer not found")

var ErrRestaurantIDMismatch = errors.New("restaurant id mismatch")
var ErrRestaurantNotFound = errors.New("restaurant not found")
var ErrRestaurantNotAcceptingOrders = errors.New("restaurant not accepting orders")
var ErrMenuItemNotFound = errors.New("menu item not found")
var ErrMenuItemNotAvailable = errors.New("item not available")
var ErrMinOrderNotReached = errors.New("min order amount not met")

var ErrDeliveryProvider = errors.New("delivery provider error")
var ErrInvalidDeliveryAddress = errors.New("invalid delivery address")

var ErrCartNotFound = errors.New("cart not found")
var ErrCartLimit = errors.New("invalid quantity")
var ErrCartEmpty = errors.New("cart empty")
var ErrInvalidInstructions = errors.New("invalid instructions")
