package entities

import "github.com/google/uuid"

type User struct {
	ID            uuid.UUID
	Email         string `encryption:"-salt"`
	VerifiedEmail bool
	Login         string `encryption:"-salt"`
	Password      string
	Picture       string `encryption:""`
}
