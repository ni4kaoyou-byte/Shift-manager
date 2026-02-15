package changerequest

import (
	"context"

	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/change_request"
)

type MemoryStore struct{}

func New() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Ping(_ context.Context) error {
	return nil
}

var _ usecase.Store = (*MemoryStore)(nil)
