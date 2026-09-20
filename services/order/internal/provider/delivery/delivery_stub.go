package delivery

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StubProvider struct {
}

func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

func (d *StubProvider) Create(_ context.Context, _ uuid.UUID) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (d *StubProvider) Start(_ context.Context, _ uuid.UUID) error {
	time.Sleep(1 * time.Second)
	return nil
}
