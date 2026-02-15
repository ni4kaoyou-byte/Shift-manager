package period

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/period"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/period"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/period"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
