package repositories

import (
	"context"

	"github.com/segmentio/kafka-go"
	kafkaModels "github.com/stdyum/api-common/modules/kafka"
)

func (r *repository) Send(ctx context.Context, message kafkaModels.Message) error {
	headers := make([]kafka.Header, 0, len(message.Headers))
	for key, value := range message.Headers {
		header := kafka.Header{
			Key:   key,
			Value: []byte(value),
		}
		headers = append(headers, header)
	}

	kafkaMessage := kafka.Message{
		Key:     []byte(message.Key),
		Value:   []byte(message.Value),
		Headers: headers,
	}

	return r.kafka.WriteMessages(ctx, kafkaMessage)
}
