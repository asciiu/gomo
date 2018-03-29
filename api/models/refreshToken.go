package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id            string
	UserId        string
	Selector      string
	Authenticator string
	TokenHash     string
	ExpiresOn     time.Time
}

func (token *RefreshToken) Update(selector, hash string, expire time.Time) {
	token.Selector = selector
	token.TokenHash = hash
	token.ExpiresOn = expire
}

func (token *RefreshToken) Renew(expire time.Time) string {
	token.ExpiresOn = expire

	// random selector
	selector := make([]byte, 16)
	rand.Read(selector)
	selectStr := base64.StdEncoding.EncodeToString(selector)
	token.Selector = selectStr

	// random authenticator
	authenticator := make([]byte, 64)
	rand.Read(authenticator)
	authenticatorStr := base64.StdEncoding.EncodeToString(authenticator)

	h := sha256.New()
	h.Write([]byte(authenticatorStr))
	token.TokenHash = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s:%s",
		base64.StdEncoding.EncodeToString(selector),
		authenticatorStr)
}

func (token *RefreshToken) Compare(authenticator string) bool {
	h := sha256.New()
	h.Write([]byte(authenticator))
	return token.TokenHash == base64.StdEncoding.EncodeToString(h.Sum(nil))
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

func NewRefreshToken(userId string) *RefreshToken {
	newId := uuid.New()

	token := RefreshToken{
		Id:     newId.String(),
		UserId: userId,
	}
	return &token
}
