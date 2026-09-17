package domain

type Restaurant struct {
	ID                string
	Name              string
	AcceptingOrders   bool
	MinimumOrderMinor int64
	Currency          string
}
