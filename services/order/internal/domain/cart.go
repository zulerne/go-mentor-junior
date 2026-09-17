package domain

type CartItem struct {
	MenuItemID     string
	Name           string
	UnitPriceMinor int64
	Currency       string
	Quantity       int32
	Instructions   string
}

type Cart struct {
	RestaurantID  *string
	Items         []CartItem
	SubtotalMinor int64
	Currency      *string
}
