package notifications

import (
	"github.com/segmentio/kafka-go"
	"github.com/stdyum/api-auth/internal/modules/notifications/controllers"
	"github.com/stdyum/api-auth/internal/modules/notifications/repositories"
)

type Notifications struct {
	controllers.Controller
}

func NewNotifications(kafka *kafka.Writer) (Notifications, error) {
	repo := repositories.NewRepository(kafka)
	ctrl := controllers.NewController(repo)

	return Notifications{Controller: ctrl}, nil
}
