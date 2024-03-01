package controllers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stdyum/api-auth/internal/app/dto"
	"github.com/stdyum/api-auth/internal/app/entities"
	"github.com/stdyum/api-auth/internal/app/models"
	"github.com/stdyum/api-auth/pkg/hash"
	jwt "github.com/stdyum/api-auth/pkg/jwt/entities"
	"github.com/stdyum/api-common/modules/notifications"
	"github.com/stdyum/api-common/proto/impl/auth"
)

func (c *controller) CreateJWTClaims(ctx context.Context, id string, userIDRaw string) (models.Claims, error) {
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		return models.Claims{}, err
	}

	user, err := c.repository.GetUserByID(ctx, userID)
	if err != nil {
		return models.Claims{}, err
	}

	if err = c.encryption.Decrypt(&user); err != nil {
		return models.Claims{}, err
	}

	claims := models.Claims{
		IDClaims: jwt.IDClaims{
			ID: id,
		},
		UserId:        user.ID,
		Login:         user.Login,
		PictureURL:    user.Picture,
		Email:         user.Email,
		VerifiedEmail: user.VerifiedEmail,
	}
	return claims, nil
}

func (c *controller) SignUp(ctx context.Context, requestDTO dto.SignUpRequestDTO) (dto.SignUpResponseDTO, error) {
	password, err := hash.Hash(requestDTO.Password)
	if err != nil {
		return dto.SignUpResponseDTO{}, err
	}

	user := entities.User{
		ID:            uuid.New(),
		Email:         requestDTO.Email,
		VerifiedEmail: false,
		Login:         requestDTO.Login,
		Password:      password,
		Picture:       requestDTO.Picture,
	}

	encryptedUser := user
	if err = c.encryption.Encrypt(&encryptedUser); err != nil {
		return dto.SignUpResponseDTO{}, err
	}

	if err = c.repository.CreateUser(ctx, encryptedUser); err != nil {
		return dto.SignUpResponseDTO{}, err
	}

	tokenPair, err := c.jwt.Create(ctx, "0.0.0.0", user.ID.String())
	if err != nil {
		return dto.SignUpResponseDTO{}, err
	}

	go func() {
		code, err := c.confirmationCodes.CreateAndStoreCode(ctx, user.ID)
		if err != nil {
			logrus.Errorf("error creating code: %v", err)
			return
		}

		notification := notifications.Notification{
			TemplateID: "sign_up_email_verification",
			Email:      user.Email,
			Data: map[string]string{
				"login": user.Login,
				"code":  code,
			},
		}

		if err := c.notifications.Send(ctx, notification); err != nil {
			logrus.Errorf("error sending notificaion: %v", err)
			return
		}
	}()

	responseDTO := dto.SignUpResponseDTO{
		User: dto.ResponseUserDTO{
			ID:      user.ID,
			Email:   user.Email,
			Login:   user.Login,
			Picture: user.Picture,
		},
		Tokens: tokenPair,
	}

	return responseDTO, nil
}

func (c *controller) Login(ctx context.Context, requestDTO dto.LoginRequestDTO) (dto.LoginResponseDTO, error) {
	login, err := c.encryption.EncryptString(requestDTO.Login, false)
	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	user, err := c.repository.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.LoginResponseDTO{}, fmt.Errorf("%s: %w", "incorrect login or password", ErrValidation)
		}
		return dto.LoginResponseDTO{}, err
	}

	if err = c.encryption.Decrypt(&user); err != nil {
		return dto.LoginResponseDTO{}, err
	}

	if ok := hash.CompareHashAndPassword(user.Password, requestDTO.Password); !ok {
		return dto.LoginResponseDTO{}, fmt.Errorf("%s: %w", "incorrect login or password", ErrValidation)
	}

	tokenPair, err := c.jwt.Create(ctx, "0.0.0.0", user.ID.String())
	if err != nil {
		return dto.LoginResponseDTO{}, err
	}

	responseDTO := dto.LoginResponseDTO{
		User: dto.ResponseUserDTO{
			ID:      user.ID,
			Email:   user.Email,
			Login:   user.Login,
			Picture: user.Picture,
		},
		Tokens: tokenPair,
	}

	return responseDTO, nil
}

func (c *controller) ConfirmEmailByCode(ctx context.Context, requestDTO dto.ConfirmEmailByCodeRequestDTO) error {
	userId, err := c.confirmationCodes.GetUserIdByCode(ctx, requestDTO.Code)
	if err != nil {
		return err
	}

	if err = c.repository.SetEmailConfirmed(ctx, userId); err != nil {
		return err
	}

	if err = c.confirmationCodes.DeleteCode(ctx, requestDTO.Code); err != nil {
		return err
	}

	return nil
}

func (c *controller) ResetPasswordRequest(ctx context.Context, requestDTO dto.ResetPasswordRequestDTO) error {
	login, err := c.encryption.EncryptString(requestDTO.Login, false)
	if err != nil {
		return err
	}

	email, err := c.encryption.EncryptString(requestDTO.Email, false)
	if err != nil {
		return err
	}

	user, err := c.repository.GetUserByLoginAndEmail(ctx, login, email)
	if err != nil {
		return err
	}

	go func() {
		code, err := c.resetCodes.CreateAndStoreCode(ctx, user.ID)
		if err != nil {
			logrus.Errorf("error creating code: %v", err)
			return
		}

		notification := notifications.Notification{
			TemplateID: "password_reset",
			Email:      user.Email,
			Data: map[string]string{
				"login": user.Login,
				"code":  code,
			},
		}

		if err := c.notifications.Send(ctx, notification); err != nil {
			logrus.Errorf("error sending notificaion: %v", err)
			return
		}
	}()

	return nil
}

func (c *controller) ResetPasswordByCode(ctx context.Context, requestDTO dto.ResetPasswordByCodeRequestDTO) error {
	userId, err := c.resetCodes.GetUserIdByCode(ctx, requestDTO.Code)
	if err != nil {
		return err
	}

	password, err := hash.Hash(requestDTO.Password)
	if err != nil {
		return err
	}

	if err = c.repository.SetPassword(ctx, userId, password); err != nil {
		return err
	}

	if err = c.resetCodes.DeleteCode(ctx, requestDTO.Code); err != nil {
		return err
	}

	return nil
}

func (c *controller) UpdateToken(ctx context.Context, request dto.UpdateTokenRequestDTO) (dto.UpdateTokenResponseDTO, error) {
	tokens, err := c.jwt.UpdateTokensByRefresh(ctx, request.Token, "0.0.0.0")
	if err != nil {
		return dto.UpdateTokenResponseDTO{}, err
	}

	return dto.UpdateTokenResponseDTO{
		Tokens: tokens,
	}, nil
}

func (c *controller) Auth(ctx context.Context, token string) (*auth.User, error) {
	claims, err := c.jwt.ValidateAccess(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err.Error(), ErrUnauthorized)
	}

	user, err := c.repository.GetUserByID(ctx, claims.Claims.UserId)
	if err != nil {
		return nil, err
	}

	if err = c.encryption.Decrypt(&user); err != nil {
		return nil, err
	}

	userRpc := auth.User{
		Id:            user.ID.String(),
		Login:         user.Login,
		PictureUrl:    user.Picture,
		Email:         user.Email,
		VerifiedEmail: user.VerifiedEmail,
	}

	return &userRpc, nil
}
