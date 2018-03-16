package models

import (
	"github.com/satori/go.uuid"
)

type User struct {
	Id            uuid.UUID
	First         string
	Last          string
	Email         string
	EmailVerified bool
	PasswordHash  string
}
