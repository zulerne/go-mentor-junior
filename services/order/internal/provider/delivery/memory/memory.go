package memory

import (
	"context"
	"time"
)

type Delivery struct {
}

func NewDelivery() *Delivery {
	return &Delivery{}
}

func (d *Delivery) Create(ctx context.Context, order_id string) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (d *Delivery) Start(ctx context.Context, order_id string) error {
	time.Sleep(1 * time.Second)
	return nil
}
