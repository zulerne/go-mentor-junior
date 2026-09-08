package domain

type CartItem struct {
	menuItemID     string
	name           string
	unitPriceMinor int64
	currency       string
	quantity       int32
	instructions   string
}

type Cart struct {
	restaurantID  *string
	items         []CartItem
	subtotalMinor int64
	currency      *string
}
