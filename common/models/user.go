package models

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(email string, password string) *User {
	newId, _ := uuid.NewV1()

	user := User{
		Id:            newId.String(),
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
