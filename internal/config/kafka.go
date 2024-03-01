package config

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type NotificationsKafkaConfig struct {
	URL      string `env:"URL"`
	Topic    string `env:"TOPIC"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}

func NotificationsKafka(config NotificationsKafkaConfig) (*kafka.Writer, error) {
	return &kafka.Writer{
		Addr:     kafka.TCP(config.URL),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: config.User,
				Password: config.Password,
			},
		},
	}, nil
}
