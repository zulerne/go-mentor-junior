package handler

type CartItem struct {
	MenuItemID     string `json:"menu_item_id"`
	Name           string `json:"name"`
	UnitPriceMinor int64  `json:"unit_price_minor"`
	Currency       string `json:"currency"`
	Quantity       int32  `json:"quantity"`
	Instructions   string `json:"instructions"`
}

type CartResponse struct {
	RestaurantID  string     `json:"restaurant_id,omitempty"`
	Items         []CartItem `json:"items"`
	SubtotalMinor int64      `json:"subtotal_minor"`
	Currency      string     `json:"currency,omitempty"`
}

type OrderItem struct {
	MenuItemID     string `json:"menu_item_id"`
	Name           string `json:"name"`
	UnitPriceMinor int64  `json:"unit_price_minor"`
	Quantity       int32  `json:"quantity"`
	Instructions   string `json:"instructions"`
}

type OrderResponse struct {
	ID              string      `json:"id"`
	CustomerID      string      `json:"customer_id"`
	RestaurantID    string      `json:"restaurant_id"`
	Status          string      `json:"status"`
	Items           []OrderItem `json:"items"`
	SubtotalMinor   int64       `json:"subtotal_minor"`
	Currency        string      `json:"currency"`
	DeliveryAddress string      `json:"delivery_address"`
	RejectionReason string      `json:"rejection_reason,omitempty"`
	DeliveryStatus  string      `json:"delivery_status,omitempty"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}

type AllOrdersResponse struct {
	Orders []OrderResponse `json:"orders"`
}
