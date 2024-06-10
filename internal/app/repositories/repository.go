package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/stdyum/api-auth/internal/app/entities"
	"github.com/stdyum/api-auth/internal/config"
	"golang.org/x/oauth2"
)

type Repository interface {
	//AuthViaOAuth2 todo move to microservice
	AuthViaOAuth2(ctx context.Context, provider string) (string, error)
	GetUserDataFromOAuth2Token(ctx context.Context, provider string, token string) (user entities.OAuth2User, err error)

	CreateUser(ctx context.Context, user entities.User) error

	GetUserByID(ctx context.Context, id uuid.UUID) (entities.User, error)
	GetUserByLogin(ctx context.Context, login string) (entities.User, error)
	GetUserByLoginAndEmail(ctx context.Context, login string, email string) (entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)

	SetEmailConfirmed(ctx context.Context, userId uuid.UUID) error
	SetPassword(ctx context.Context, userId uuid.UUID, password string) error
}

type repository struct {
	database *sql.DB
	oauth    map[string]oauth2.Config
}

func New(database *sql.DB, oauth config.OAuthConfig) Repository {
	return &repository{
		database: database,
		oauth:    buildOAuthClients(oauth),
	}
}

func buildOAuthClients(oauth config.OAuthConfig) map[string]oauth2.Config {
	return map[string]oauth2.Config{
		"google": buildGoogleOAuthClient(oauth.Google),
	}
}

func buildGoogleOAuthClient(google config.OAuthGoogleConfig) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  google.RedirectURL,
		ClientID:     google.ClientId,
		ClientSecret: google.ClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:       "https://accounts.google.com/o/oauth2/auth",
			TokenURL:      "https://oauth2.googleapis.com/token",
			DeviceAuthURL: "https://oauth2.googleapis.com/device/code",
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}
}
