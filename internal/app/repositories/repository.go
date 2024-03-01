package repositories

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stdyum/api-auth/internal/app/entities"
)

type Repository interface {
	CreateUser(ctx context.Context, user entities.User) error

	GetUserByID(ctx context.Context, id uuid.UUID) (entities.User, error)
	GetUserByLogin(ctx context.Context, login string) (entities.User, error)
	GetUserByLoginAndEmail(ctx context.Context, login string, email string) (entities.User, error)

	SetEmailConfirmed(ctx context.Context, userId uuid.UUID) error
	SetPassword(ctx context.Context, userId uuid.UUID, password string) error
}

type repository struct {
	database *sql.DB
}

func New(database *sql.DB) Repository {
	return &repository{
		database: database,
	}
}
