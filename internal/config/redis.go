package config

import (
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	URL      string `env:"URL"`
	Password string `env:"PASSWORD"`
}

func ConnectToRedis(config RedisConfig, database int) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       database,
	}), nil
}
