package codes

import (
	"github.com/redis/go-redis/v9"
	"github.com/stdyum/api-auth/internal/modules/codes/controllers"
	"github.com/stdyum/api-auth/internal/modules/codes/repositories"
)

type Codes struct {
	controllers.Controller
}

func NewCodes(database *redis.Client) (Codes, error) {
	repo := repositories.NewRepository(database)
	ctrl := controllers.NewController(repo)

	return Codes{Controller: ctrl}, nil
}
