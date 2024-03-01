package controllers

import (
	"context"

	"github.com/stdyum/api-auth/internal/modules/notifications/repositories"
	"github.com/stdyum/api-common/modules/notifications"
)

type Controller interface {
	Send(ctx context.Context, message notifications.Notification) error
}

type controller struct {
	repository repositories.Repository
}

func NewController(repository repositories.Repository) Controller {
	return &controller{repository: repository}
}
