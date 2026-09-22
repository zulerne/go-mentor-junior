package domain

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrRestaurantNotFound = errors.New("restaurant not found")
var ErrCustomerNotFound = errors.New("customer not found")
var ErrRestaurantIDMismatch = errors.New("restaurant id mismatch")
var ErrInvalidOrderStatus = errors.New("invalid order status")
var ErrDeliveryProvider = errors.New("delivery provider error")
var ErrCartNotFound = errors.New("cart not found")
var ErrMenuItemNotFound = errors.New("menu item not found")
var ErrItemNotAvailable = errors.New("item not available")
var ErrCartLimit = errors.New("invalid quantity")
var ErrInvalidInstructions = errors.New("invalid instructions")
