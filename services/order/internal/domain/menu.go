package domain

type MenuItem struct {
	ID          string
	Name        string
	Description string
	PriceMinor  int64
	Currency    string
	Available   bool
}
