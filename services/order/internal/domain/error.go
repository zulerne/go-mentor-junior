package domain

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrRestaurantNotFound = errors.New("restaurant not found")
var ErrCustomerNotFound = errors.New("customer not found")
var ErrRestaurantIDMismatch = errors.New("restaurant id mismatch")
var ErrInvalidOrderStatus = errors.New("invalid order status")
var ErrDeliveryProvider = errors.New("delivery provider error")
