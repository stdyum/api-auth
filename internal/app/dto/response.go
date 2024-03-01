package dto

import (
	"github.com/google/uuid"
	"github.com/stdyum/api-auth/pkg/jwt/entities"
)

type ResponseUserDTO struct {
	ID      uuid.UUID `json:"id"`
	Email   string    `json:"email"`
	Login   string    `json:"login"`
	Picture string    `json:"picture"`
}

type SignUpResponseDTO struct {
	User   ResponseUserDTO    `json:"user"`
	Tokens entities.TokenPair `json:"tokens"`
}

type LoginResponseDTO struct {
	User   ResponseUserDTO    `json:"user"`
	Tokens entities.TokenPair `json:"tokens"`
}

type UpdateTokenResponseDTO struct {
	Tokens entities.TokenPair `json:"tokens"`
}
