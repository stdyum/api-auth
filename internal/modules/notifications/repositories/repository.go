package repositories

import (
	"context"

	"github.com/segmentio/kafka-go"
	kafkaModels "github.com/stdyum/api-common/modules/kafka"
)

type Repository interface {
	Send(ctx context.Context, message kafkaModels.Message) error
}

type repository struct {
	kafka *kafka.Writer
}

func NewRepository(kafka *kafka.Writer) Repository {
	return &repository{kafka: kafka}
}
