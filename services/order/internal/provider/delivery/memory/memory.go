package memory

import (
	"context"
	"time"
)

type DeliveryProvider struct {
}

func NewDeliveryProvider() *DeliveryProvider {
	return &DeliveryProvider{}
}

func (d *DeliveryProvider) Create(ctx context.Context, order_id string) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (d *DeliveryProvider) Start(ctx context.Context, order_id string) error {
	time.Sleep(1 * time.Second)
	return nil
}
