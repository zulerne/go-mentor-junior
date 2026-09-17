package delivery

import (
	"context"
	"time"
)

type StubProvider struct {
}

func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

func (d *StubProvider) Create(_ context.Context, _ string) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (d *StubProvider) Start(_ context.Context, _ string) error {
	time.Sleep(1 * time.Second)
	return nil
}
