package domain

import "github.com/google/uuid"

type Restaurant struct {
	ID                uuid.UUID
	Name              string
	AcceptingOrders   bool
	MinimumOrderMinor int64
	Currency          string
}
