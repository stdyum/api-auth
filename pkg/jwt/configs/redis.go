package configs

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	CronPattern   string
	Expire        time.Duration
	RefreshExpire time.Duration
	Timeout       time.Duration
	Secret        string
	Client        *redis.Client
}
