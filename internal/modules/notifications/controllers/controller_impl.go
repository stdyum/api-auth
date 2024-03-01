package controllers

import (
	"context"
	"encoding/json"

	"github.com/stdyum/api-common/modules/kafka"
	"github.com/stdyum/api-common/modules/notifications"
)

func (r *controller) Send(ctx context.Context, notification notifications.Notification) error {
	bytes, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	message := kafka.Message{
		Key:   "email",
		Value: string(bytes),
	}
	return r.repository.Send(ctx, message)
}
