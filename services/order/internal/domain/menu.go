package domain

type MenuItem struct {
	id          string
	name        string
	description string
	priceMinor  int64
	currency    string
	available   bool
}
