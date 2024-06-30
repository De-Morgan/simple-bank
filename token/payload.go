package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Payload struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (payload *Payload, err error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return
	}
	payload = &Payload{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	return
}

var ErrExpiredToken = errors.New("token has expired")

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt.Time) {
		return ErrExpiredToken
	}
	return nil
}
