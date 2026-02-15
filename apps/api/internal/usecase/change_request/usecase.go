package changerequest

import (
	"context"

	repository "github.com/ni4kaoyou-byte/shift-manager/apps/api/internal/repository/change_request"
)

type Usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Ping(ctx context.Context) error {
	return u.repo.Ping(ctx)
}
