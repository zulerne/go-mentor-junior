package domain

import "time"

type OrderStatus int

const (
	Pending OrderStatus = iota
	Accepted
	Preparing
	ReadyForPickup
	Rejected
	Cancelled
)

type OrderDeliveryStatus int

const (
	WaitingForPreparation OrderDeliveryStatus = iota
	InDelivery
)

type OrderItem struct {
	menuItemID     string
	name           string
	unitPriceMinor int64
	quantity       int32
	instructions   string
}

type Order struct {
	id              string
	customerID      string
	restaurantID    string
	status          OrderStatus
	items           []OrderItem
	subtotalMinor   int64
	currency        string
	deliveryAddress string
	rejectionReason *string
	deliveryStatus  *OrderDeliveryStatus
	createdAt       time.Time
	updatedAt       time.Time
}
