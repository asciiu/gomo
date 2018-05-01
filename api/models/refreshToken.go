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
	ID            string
	UserID        string
	Selector      string
	Authenticator string
	TokenHash     string
	ExpiresOn     time.Time
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

// Returns true if the authenticator matches the tokens hash and
// the token has not expired
func (token *RefreshToken) Valid(authenticator string) bool {
	h := sha256.New()
	h.Write([]byte(authenticator))
	return token.TokenHash == base64.StdEncoding.EncodeToString(h.Sum(nil)) &&
		token.ExpiresOn.After(time.Now())
}

func NewRefreshToken(userID string) *RefreshToken {
	newID := uuid.New()

	token := RefreshToken{
		ID:     newID.String(),
		UserID: userID,
	}
	return &token
}
