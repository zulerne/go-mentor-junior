package delivery

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type StubProvider struct {
	mu sync.Mutex
}

func NewStubProvider() *StubProvider {
	return &StubProvider{}
}

func (d *StubProvider) Create(_ context.Context, _ uuid.UUID) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	time.Sleep(1 * time.Second)
	return nil
}

func (d *StubProvider) Start(_ context.Context, _ uuid.UUID) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	time.Sleep(1 * time.Second)
	return nil
}
