package internal

import (
	"github.com/stdyum/api-auth/internal/app"
	"github.com/stdyum/api-auth/internal/config"
	"github.com/stdyum/api-auth/internal/modules/codes"
	"github.com/stdyum/api-auth/internal/modules/notifications"
	"github.com/stdyum/api-auth/pkg/encryption"
)

func App() error {
	db, err := config.ConnectToDatabase(config.Config.Database)
	if err != nil {
		return err
	}

	redisJWT, err := config.ConnectToRedis(config.Config.Redis, 0)
	if err != nil {
		return err
	}

	jwt, err := config.JWTWithRedis(config.Config.JWT, redisJWT)
	if err != nil {
		return err
	}

	encrypt, err := encryption.NewEncryption(config.Config.EncryptionSecret)
	if err != nil {
		return err
	}

	redisConfirmationCodes, err := config.ConnectToRedis(config.Config.Redis, 1)
	if err != nil {
		return err
	}

	redisResetCodes, err := config.ConnectToRedis(config.Config.Redis, 2)
	if err != nil {
		return err
	}

	notificationsKafka, err := config.NotificationsKafka(config.Config.NotificationsKafka)
	if err != nil {
		return err
	}

	notificationsModule, err := notifications.NewNotifications(notificationsKafka)
	if err != nil {
		return err
	}

	confirmationCodesModule, err := codes.NewCodes(redisConfirmationCodes)
	if err != nil {
		return err
	}

	resetCodesModule, err := codes.NewCodes(redisResetCodes)
	if err != nil {
		return err
	}

	routes, err := app.New(jwt, encrypt, notificationsModule, confirmationCodesModule, resetCodesModule, db)
	if err != nil {
		return err
	}

	routes.Ports = config.Config.Ports

	return routes.Run()
}
