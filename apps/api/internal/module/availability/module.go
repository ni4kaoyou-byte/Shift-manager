package availability

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/availability"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/availability"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/availability"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
