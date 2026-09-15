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

func (d *DeliveryProvider) Create(_ context.Context, _ string) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (d *DeliveryProvider) Start(_ context.Context, _ string) error {
	time.Sleep(1 * time.Second)
	return nil
}
