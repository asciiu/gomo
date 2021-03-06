package models

import (
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(first, last, email, password string) *User {
	newID := uuid.New()

	user := User{
		ID:            newID.String(),
		First:         first,
		Last:          last,
		Email:         email,
		EmailVerified: false,
		PasswordHash:  HashAndSalt([]byte(password)),
	}
	return &user
}

func HashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	// hash this using a server secret key
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

type User struct {
	ID            string
	First         string
	Last          string
	Email         string
	EmailVerified bool
	PasswordHash  string
	Salt          string
}

type UserInfo struct {
	UserID string `json:"userID"`
	First  string `json:"first"`
	Last   string `json:"last"`
	Email  string `json:"email"`
}

func (user *User) Info() *UserInfo {
	return &UserInfo{
		UserID: user.ID,
		First:  user.First,
		Last:   user.Last,
		Email:  user.Email,
	}
}
