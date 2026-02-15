package assignment

import "context"

type Store interface {
	Ping(ctx context.Context) error
}

type Service struct {
	store Store
}

func New(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Ping(ctx context.Context) error {
	return s.store.Ping(ctx)
}
