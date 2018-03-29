package models

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id        string
	UserId    string
	Selector  string
	TokenHash string
	ExpiresOn time.Time
}

func NewSelectorAuth() string {
	// random selector
	selector := make([]byte, 16)
	rand.Read(selector)

	// random authenticator
	authenticator := make([]byte, 64)
	rand.Read(authenticator)

	return fmt.Sprintf("%s:%s",
		base64.StdEncoding.EncodeToString(selector),
		base64.StdEncoding.EncodeToString(authenticator))
}

func NewRefreshToken(userId, selectorAuth string, expiresOn time.Time) *RefreshToken {
	newId := uuid.New()
	pts := strings.Split(selectorAuth, ":")

	token := RefreshToken{
		Id:        newId.String(),
		UserId:    userId,
		Selector:  pts[0],
		TokenHash: pts[1],
		ExpiresOn: expiresOn,
	}
	return &token
}
