package auditlog

import (
	"context"

	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/audit_log"
)

type MemoryStore struct{}

func New() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Ping(_ context.Context) error {
	return nil
}

var _ usecase.Store = (*MemoryStore)(nil)
