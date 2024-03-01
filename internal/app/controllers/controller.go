package controllers

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stdyum/api-auth/internal/app/dto"
	"github.com/stdyum/api-auth/internal/app/models"
	"github.com/stdyum/api-auth/internal/app/repositories"
	"github.com/stdyum/api-auth/internal/modules/codes"
	"github.com/stdyum/api-auth/internal/modules/notifications"
	"github.com/stdyum/api-auth/pkg/encryption"
	jwt "github.com/stdyum/api-auth/pkg/jwt/controllers"
	"github.com/stdyum/api-common/proto/impl/auth"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrValidation   = errors.New("validation")
)

type Controller interface {
	CreateJWTClaims(ctx context.Context, id string, userID string) (models.Claims, error)

	SignUp(ctx context.Context, requestDTO dto.SignUpRequestDTO) (dto.SignUpResponseDTO, error)
	Login(ctx context.Context, requestDTO dto.LoginRequestDTO) (dto.LoginResponseDTO, error)

	ConfirmEmailByCode(ctx context.Context, requestDTO dto.ConfirmEmailByCodeRequestDTO) error

	ResetPasswordRequest(ctx context.Context, requestDTO dto.ResetPasswordRequestDTO) error
	ResetPasswordByCode(ctx context.Context, requestDTO dto.ResetPasswordByCodeRequestDTO) error

	UpdateToken(ctx context.Context, request dto.UpdateTokenRequestDTO) (dto.UpdateTokenResponseDTO, error)

	Auth(ctx context.Context, token string) (*auth.User, error)
}

type controller struct {
	encryption        encryption.Encryption
	jwt               jwt.Controller[models.Claims]
	notifications     notifications.Notifications
	confirmationCodes codes.Codes
	resetCodes        codes.Codes

	repository repositories.Repository
}

func New(jwt jwt.Controller[models.Claims], encryption encryption.Encryption, notifications notifications.Notifications, confirmationCodes codes.Codes, resetCodes codes.Codes, repository repositories.Repository) Controller {
	return &controller{
		encryption:        encryption,
		jwt:               jwt,
		notifications:     notifications,
		confirmationCodes: confirmationCodes,
		resetCodes:        resetCodes,
		repository:        repository,
	}
}
