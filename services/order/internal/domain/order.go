package domain

import "time"

type OrderStatus string

const (
	Pending        OrderStatus = "PENDING"
	Accepted       OrderStatus = "ACCEPTED"
	Preparing      OrderStatus = "PREPARING"
	ReadyForPickup OrderStatus = "READY_FOR_PICKUP"
	Rejected       OrderStatus = "REJECTED"
	Cancelled      OrderStatus = "CANCELLED"
)

type OrderDeliveryStatus string

const (
	WaitingForPreparation OrderDeliveryStatus = "WAITING_FOR_PREPARATION"
	InDelivery            OrderDeliveryStatus = "IN_DELIVERY"
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
