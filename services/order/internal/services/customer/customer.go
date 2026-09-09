package customer

import "log/slog"

type Store interface {
}

type Service struct {
	store Store
	log   *slog.Logger
}

func New(store Store, log *slog.Logger) *Service {
	log = log.With("component", "customer")
	c := &Service{
		store: store,
		log:   log,
	}

	return c
}
