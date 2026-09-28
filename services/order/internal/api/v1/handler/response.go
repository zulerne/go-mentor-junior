package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/zulerne/go-mentor-junior/order/internal/domain"
)

type CartItem struct {
	MenuItemID     uuid.UUID `json:"menu_item_id"`
	Name           string    `json:"name"`
	UnitPriceMinor int64     `json:"unit_price_minor"`
	Currency       string    `json:"currency"`
	Quantity       int32     `json:"quantity"`
	Instructions   string    `json:"instructions"`
}

type CartResponse struct {
	RestaurantID  *uuid.UUID `json:"restaurant_id"`
	Items         []CartItem `json:"items"`
	SubtotalMinor int64      `json:"subtotal_minor"`
	Currency      string     `json:"currency,omitempty"`
}

type OrderItem struct {
	MenuItemID     uuid.UUID `json:"menu_item_id"`
	Name           string    `json:"name"`
	UnitPriceMinor int64     `json:"unit_price_minor"`
	Quantity       int32     `json:"quantity"`
	Instructions   string    `json:"instructions"`
}

type OrderResponse struct {
	ID              uuid.UUID   `json:"id"`
	CustomerID      uuid.UUID   `json:"customer_id"`
	RestaurantID    uuid.UUID   `json:"restaurant_id"`
	Status          string      `json:"status"`
	Items           []OrderItem `json:"items"`
	SubtotalMinor   int64       `json:"subtotal_minor"`
	Currency        string      `json:"currency"`
	DeliveryAddress string      `json:"delivery_address"`
	RejectionReason string      `json:"rejection_reason,omitempty"`
	DeliveryStatus  string      `json:"delivery_status,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

type AllOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}

func cartItemToDTO(src domain.CartItem) CartItem {
	return CartItem{
		MenuItemID:     src.MenuItemID,
		Name:           src.Name,
		UnitPriceMinor: src.UnitPriceMinor,
		Currency:       src.Currency,
		Quantity:       src.Quantity,
		Instructions:   src.Instructions,
	}
}

func cartToDTO(src domain.Cart) CartResponse {
	items := make([]CartItem, 0, len(src.Items))
	for _, item := range src.Items {
		items = append(items, cartItemToDTO(item))
	}
	var restaurantID *uuid.UUID
	if src.RestaurantID != uuid.Nil {
		id := src.RestaurantID
		restaurantID = &id
	}
	return CartResponse{
		RestaurantID:  restaurantID,
		Items:         items,
		SubtotalMinor: src.SubtotalMinor,
		Currency:      src.Currency,
	}
}

func orderItemToDTO(src domain.OrderItem) OrderItem {
	return OrderItem{
		MenuItemID:     src.MenuItemID,
		Name:           src.Name,
		UnitPriceMinor: src.UnitPriceMinor,
		Quantity:       src.Quantity,
		Instructions:   src.Instructions,
	}
}

func orderToDTO(src domain.Order) OrderResponse {
	items := make([]OrderItem, 0, len(src.Items))
	for _, item := range src.Items {
		items = append(items, orderItemToDTO(item))
	}
	return OrderResponse{
		ID:              src.ID,
		CustomerID:      src.CustomerID,
		RestaurantID:    src.RestaurantID,
		Status:          string(src.Status),
		Items:           items,
		SubtotalMinor:   src.SubtotalMinor,
		Currency:        src.Currency,
		DeliveryAddress: src.DeliveryAddress,
		RejectionReason: src.RejectionReason,
		DeliveryStatus:  string(src.DeliveryStatus),
		CreatedAt:       src.CreatedAt,
		UpdatedAt:       src.UpdatedAt,
	}
}

func ordersToDTOs(src []domain.Order) []OrderResponse {
	result := make([]OrderResponse, 0, len(src))
	for _, order := range src {
		result = append(result, orderToDTO(order))
	}
	return result
}
