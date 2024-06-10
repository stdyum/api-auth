package dto

import (
	"time"
)

type AuthViaOAuth2Request struct {
	Provider string `json:"provider"`
}

type AuthViaOAuth2CallbackRequest struct {
	Provider string `json:"provider"`
	Code     string `json:"code"`
}

type SignUpRequestDTO struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Picture  string `json:"picture"`
}

type LoginRequestDTO struct {
	Login               string    `json:"login"`
	Password            string    `json:"password"`
	SessionExpirationAt time.Time `json:"sessionExpirationAt"`
}

type ConfirmEmailByCodeRequestDTO struct {
	Code string `json:"code"`
}

type ResetPasswordRequestDTO struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

type ResetPasswordByCodeRequestDTO struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}

type UpdateTokenRequestDTO struct {
	Token string `json:"refresh"`
}
