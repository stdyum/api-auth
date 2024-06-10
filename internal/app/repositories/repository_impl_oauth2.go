package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/stdyum/api-auth/internal/app/entities"
	"golang.org/x/oauth2"
)

func (r *repository) AuthViaOAuth2(ctx context.Context, provider string) (string, error) {
	config, err := r.getConfigByProviderName(ctx, provider)
	if err != nil {
		return "", err
	}

	return config.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

func (r *repository) GetUserDataFromOAuth2Token(ctx context.Context, provider string, code string) (user entities.OAuth2User, err error) {
	config, err := r.getConfigByProviderName(ctx, provider)
	if err != nil {
		return entities.OAuth2User{}, err
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return entities.OAuth2User{}, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	// todo url according to provider
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return entities.OAuth2User{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&user)
	return
}

func (r *repository) getConfigByProviderName(_ context.Context, provider string) (oauth2.Config, error) {
	config, ok := r.oauth[provider]
	if !ok {
		return oauth2.Config{}, errors.New("invalid provider")
	}

	return config, nil
}
