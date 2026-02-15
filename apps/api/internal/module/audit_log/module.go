package auditlog

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/audit_log"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/audit_log"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/audit_log"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
