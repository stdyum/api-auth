package config

type OAuthGoogleConfig struct {
	RedirectURL  string `env:"REDIRECT_URL"`
	ClientId     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
}
