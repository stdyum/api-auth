package controllers

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stdyum/api-auth/internal/modules/codes/repositories"
)

var (
	ErrNilUUID = errors.New("nil uuid")
)

type Controller interface {
	CreateAndStoreCode(ctx context.Context, userId uuid.UUID) (string, error)

	DeleteCode(ctx context.Context, code string) error

	GetUserIdByCode(ctx context.Context, code string) (uuid.UUID, error)

	StoreCode(ctx context.Context, userId uuid.UUID, code string) error
}

type controller struct {
	repository repositories.Repository
}

func NewController(repository repositories.Repository) Controller {
	return &controller{
		repository: repository,
	}
}
