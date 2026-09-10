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
	MenuItemId     string
	Name           string
	UnitPriceMinor int64
	Quantity       int32
	Instructions   string
}

type Order struct {
	ID              string
	CustomerID      string
	RestaurantID    string
	Status          OrderStatus
	Items           []OrderItem
	SubtotalMinor   int64
	Currency        string
	DeliveryAddress string
	RejectionReason *string
	DeliveryStatus  *OrderDeliveryStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
