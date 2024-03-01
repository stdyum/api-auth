package models

import (
	"github.com/google/uuid"
	jwt "github.com/stdyum/api-auth/pkg/jwt/controllers"
	"github.com/stdyum/api-auth/pkg/jwt/entities"
)

type JWT = jwt.Controller[Claims]

type Claims struct {
	entities.IDClaims
	UserId        uuid.UUID `json:"userID"`
	Login         string    `json:"login"`
	PictureURL    string    `json:"pictureURL"`
	Email         string    `json:"email"`
	VerifiedEmail bool      `json:"verifiedEmail"`
}
