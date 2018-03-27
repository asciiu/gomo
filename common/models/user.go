package models

import (
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(first, last, email, password string) *User {
	newId := uuid.New()

	user := User{
		Id:            newId.String(),
		First:         first,
		Last:          last,
		Email:         email,
		EmailVerified: false,
		PasswordHash:  hashAndSalt([]byte(password)),
	}
	return &user
}

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

type User struct {
	Id            string
	First         string
	Last          string
	Email         string
	EmailVerified bool
	PasswordHash  string
	Salt          string
}
