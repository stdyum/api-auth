package app

import (
	"database/sql"

	"github.com/stdyum/api-auth/internal/app/controllers"
	"github.com/stdyum/api-auth/internal/app/errors"
	"github.com/stdyum/api-auth/internal/app/handlers"
	"github.com/stdyum/api-auth/internal/app/models"
	"github.com/stdyum/api-auth/internal/app/repositories"
	"github.com/stdyum/api-auth/internal/config"
	"github.com/stdyum/api-auth/internal/modules/codes"
	"github.com/stdyum/api-auth/internal/modules/notifications"
	"github.com/stdyum/api-auth/pkg/encryption"
	"github.com/stdyum/api-common/server"
)

func New(jwt models.JWT, encrypt encryption.Encryption, oauth config.OAuthConfig, url string, notifications notifications.Notifications, confirmationCodes codes.Codes, resetCodes codes.Codes, database *sql.DB) (server.Routes, error) {
	repo := repositories.New(database, oauth)

	ctrl := controllers.New(jwt, encrypt, notifications, confirmationCodes, resetCodes, repo)
	jwt.SetCreateClaimsFunc(ctrl.CreateJWTClaims)

	errors.Register()

	httpHndl := handlers.NewHTTP(ctrl, url)
	grpcHndl := handlers.NewGRPC(ctrl)

	routes := server.Routes{
		GRPC: grpcHndl,
		HTTP: httpHndl,
	}

	return routes, nil
}
