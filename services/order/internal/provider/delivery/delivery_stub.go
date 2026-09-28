package delivery

import (
	"context"

	"github.com/google/uuid"
)

type StubProvider struct{}

func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

func (d *StubProvider) Create(_ context.Context, _ uuid.UUID) error {
	return nil
}

func (d *StubProvider) Start(_ context.Context, _ uuid.UUID) error {
	return nil
}
