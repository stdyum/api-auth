package config

import (
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/stdyum/api-auth/internal/app/models"
	"github.com/stdyum/api-auth/pkg/jwt"
	jwtConfig "github.com/stdyum/api-auth/pkg/jwt/configs"
	"time"
)

type JWTConfig struct {
	Secret        string `env:"SECRET"`
	Timeout       int    `env:"TIMEOUT"`
	Expire        int    `env:"EXPIRE"`
	RefreshExpire int    `env:"EXPIRE_REFRESH"`
}

func JWTWithRedis(config JWTConfig, client *redis.Client) (models.JWT, error) {
	jwtRedisConfig := jwtConfig.Redis{
		Expire:        time.Second * time.Duration(config.Expire),
		RefreshExpire: time.Second * time.Duration(config.RefreshExpire),
		Timeout:       time.Second * time.Duration(config.Timeout),
		Secret:        config.Secret,
		Client:        client,
	}

	return jwt.NewWithRedisCfg[models.Claims](jwtRedisConfig), nil
}
