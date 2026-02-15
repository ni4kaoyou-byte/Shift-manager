package changerequest

import (
	handler "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/handler/change_request"
	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/change_request"
	usecase "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/usecase/change_request"
)

func NewHandler() *handler.Handler {
	return handler.NewHandler(usecase.New(repository.New()))
}
