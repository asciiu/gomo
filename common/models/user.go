package models

type User struct {
	Id            string
	First         string
	Last          string
	Email         string
	EmailVerified bool
	PasswordHash  string
	Salt          string
}
