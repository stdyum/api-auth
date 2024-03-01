package repositories

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	DeleteCode(ctx context.Context, code string) error

	GetUserIdByCode(ctx context.Context, code string) (string, error)

	StoreCode(ctx context.Context, userId string, code string) error
}

type repository struct {
	database *redis.Client
}

func NewRepository(database *redis.Client) Repository {
	return &repository{
		database: database,
	}
}
