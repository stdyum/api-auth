package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stdyum/api-auth/internal/app/dto"
	"github.com/stdyum/api-auth/internal/app/entities"
	appModels "github.com/stdyum/api-auth/internal/app/models"
	"github.com/stdyum/api-auth/pkg/hash"
	jwt "github.com/stdyum/api-auth/pkg/jwt/entities"
	"github.com/stdyum/api-auth/pkg/jwt/repositories"
	"github.com/stdyum/api-common/errors"
	"github.com/stdyum/api-common/models"
	"github.com/stdyum/api-common/modules/notifications"
)

func (c *controller) CreateJWTClaims(ctx context.Context, id string, userIDRaw string) (appModels.Claims, error) {
	userID, err := uuid.Parse(userIDRaw)
	if err != nil {
		return appModels.Claims{}, err
	}

	user, err := c.repository.GetUserByID(ctx, userID)
	if err != nil {
		return appModels.Claims{}, err
	}

	if err = c.encryption.Decrypt(&user); err != nil {
		return appModels.Claims{}, err
	}

	claims := appModels.Claims{
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

func (c *controller) AuthViaOAuth2(ctx context.Context, request dto.AuthViaOAuth2Request) (string, error) {
	return c.repository.AuthViaOAuth2(ctx, request.Provider)
}

func (c *controller) AuthViaOAuth2Callback(ctx context.Context, request dto.AuthViaOAuth2CallbackRequest) (jwt.TokenPair, error) {
	oauthUser, err := c.repository.GetUserDataFromOAuth2Token(ctx, request.Provider, request.Code)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	encryptedEmail, err := c.encryption.EncryptString(oauthUser.Email, false)
	if err != nil {
		return jwt.TokenPair{}, err
	}

	user, err := c.repository.GetUserByEmail(ctx, encryptedEmail)
	if errors.Is(sql.ErrNoRows, err) {
		user = entities.User{
			ID:            uuid.New(),
			Email:         oauthUser.Email,
			VerifiedEmail: true,
			Login:         oauthUser.Name,
			Picture:       oauthUser.Picture,
		}

		encryptedUser := user
		if err = c.encryption.Encrypt(&encryptedUser); err != nil {
			return jwt.TokenPair{}, err
		}

		if err = c.repository.CreateUser(ctx, encryptedUser); err != nil {
			return jwt.TokenPair{}, err
		}
	}

	return c.jwt.Create(ctx, "0.0.0.0", user.ID.String())
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

	var tokenPair jwt.TokenPair
	if requestDTO.SessionExpirationAt.Before(time.Now()) {
		tokenPair, err = c.jwt.Create(ctx, "0.0.0.0", user.ID.String())
	} else {
		tokenPair.Access, err = c.jwt.CreateAccessWithExpireTime(ctx, "0.0.0.0", user.ID.String(), requestDTO.SessionExpirationAt)
	}
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

func (c *controller) Self(ctx context.Context, token string) (dto.ResponseUserDTO, error) {
	user, err := c.Auth(ctx, token)
	if err != nil {
		return dto.ResponseUserDTO{}, err
	}

	return dto.ResponseUserDTO{
		ID:      user.ID,
		Email:   user.Email,
		Login:   user.Login,
		Picture: user.PictureUrl,
	}, err
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
		if errors.Is(repositories.NotValidRefreshTokenErr, err) {
			return dto.UpdateTokenResponseDTO{}, errors.Wrap(ErrUnauthorized, err)
		}

		return dto.UpdateTokenResponseDTO{}, err
	}

	return dto.UpdateTokenResponseDTO{
		Tokens: tokens,
	}, nil
}

func (c *controller) Auth(ctx context.Context, token string) (models.User, error) {
	claims, err := c.jwt.ValidateAccess(ctx, token)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", err.Error(), ErrUnauthorized)
	}

	user, err := c.repository.GetUserByID(ctx, claims.Claims.UserId)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return models.User{}, errors.WrapString(ErrUnauthorized, "user does not exist")
		}

		return models.User{}, err
	}

	if err = c.encryption.Decrypt(&user); err != nil {
		return models.User{}, err
	}

	u := models.User{
		ID:            user.ID,
		Login:         user.Login,
		PictureUrl:    user.Picture,
		Email:         user.Email,
		VerifiedEmail: user.VerifiedEmail,
	}

	return u, nil
}
