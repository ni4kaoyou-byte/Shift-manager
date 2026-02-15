package assignment

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/assignment"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/assignment"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/assignment"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
