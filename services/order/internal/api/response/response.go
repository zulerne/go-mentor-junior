package response

// TODO (review): I don't sure about where these structs should be placed. And is it okay that these structs almost the same as domain structs?

type CardItem struct {
	MenuItemID     string `json:"menu_item_id"`
	Name           string `json:"name"`
	UnitPriceMinor int64  `json:"unit_price_minor"`
	Currency       string `json:"currency"`
	Quantity       int32  `json:"quantity"`
	Instructions   string `json:"instructions"`
}

type Cart struct {
	RestaurantID  *string    `json:"restaurant_id"`
	Items         []CardItem `json:"items"`
	SubtotalMinor int64      `json:"subtotal_minor"`
	Currency      *string    `json:"currency"`
}

type OrderItem struct {
	MenuItemID     string `json:"menu_item_id"`
	Name           string `json:"name"`
	UnitPriceMinor int64  `json:"unit_price_minor"`
	Quantity       int32  `json:"quantity"`
	Instructions   string `json:"instructions"`
}

type Order struct {
	ID              string      `json:"id"`
	CustomerID      string      `json:"customer_id"`
	RestaurantID    string      `json:"restaurant_id"`
	Status          string      `json:"status"`
	Items           []OrderItem `json:"items"`
	SubtotalMinor   int64       `json:"subtotal_minor"`
	Currency        string      `json:"currency"`
	DeliveryAddress string      `json:"delivery_address"`
	RejectionReason *string     `json:"rejection_reason"`
	DeliveryStatus  *string     `json:"delivery_status"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
}

type AllOrders struct {
	Orders []Order `json:"orders"`
}
