package config

type OAuthConfig struct {
	Google OAuthGoogleConfig `env:"GOOGLE"`
}
