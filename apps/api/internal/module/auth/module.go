package auth

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/auth"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/auth"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/auth"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
