package domain

import "github.com/google/uuid"

type CartItem struct {
	MenuItemID     uuid.UUID
	Name           string
	UnitPriceMinor int64
	Currency       string
	Quantity       int32
	Instructions   string
}

type Cart struct {
	RestaurantID  uuid.UUID
	Items         []CartItem
	SubtotalMinor int64
	Currency      string
}
