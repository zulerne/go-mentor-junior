package domain

import (
	"time"

	"github.com/google/uuid"
)

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
	UnspecifiedOrderDeliveryStatus OrderDeliveryStatus = ""
	WaitingForPreparation          OrderDeliveryStatus = "WAITING_FOR_PREPARATION"
	InDelivery                     OrderDeliveryStatus = "IN_DELIVERY"
)

// TODO: Add micro/mini types for id fields

type OrderItem struct {
	MenuItemID     uuid.UUID
	Name           string
	UnitPriceMinor int64
	Quantity       int32
	Instructions   string
}

type Order struct {
	ID              uuid.UUID
	CustomerID      uuid.UUID
	RestaurantID    uuid.UUID
	Status          OrderStatus
	Items           []OrderItem
	SubtotalMinor   int64
	Currency        string
	DeliveryAddress string
	RejectionReason string
	DeliveryStatus  OrderDeliveryStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
