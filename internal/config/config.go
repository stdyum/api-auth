package config

import (
	"github.com/stdyum/api-common/env"
	"github.com/stdyum/api-common/server"
)

var Config Model

type Model struct {
	Ports              server.PortConfig        `env:"PORT"`
	EncryptionSecret   string                   `env:"ENCRYPTION_SECRET"`
	Database           DatabaseConfig           `env:"DATABASE"`
	Redis              RedisConfig              `env:"REDIS"`
	NotificationsKafka NotificationsKafkaConfig `env:"KAFKA_NOTIFICATIONS"`
	JWT                JWTConfig                `env:"JWT"`
}

func init() {
	err := env.Fill(&Config)
	if err != nil {
		panic("cannot fill config: " + err.Error())
	}
}
