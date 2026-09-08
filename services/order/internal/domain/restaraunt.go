package domain

type Restaurant struct {
	id                string
	name              string
	acceptingOrders   bool
	minimumOrderMinor int64
	currency          string
}
