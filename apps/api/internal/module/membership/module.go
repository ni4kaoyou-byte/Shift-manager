package membership

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/membership"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/membership"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/membership"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
