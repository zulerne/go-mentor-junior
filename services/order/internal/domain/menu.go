package domain

type MenuItem struct {
	Id          string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	Available   bool
}
