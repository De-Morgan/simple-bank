package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const minSecretKeySize = 32

type Payload struct {
	Username  string    `json:"username"`
	ID        uuid.UUID `json:"id"`
	IssuedAt  time.Time `json:"issued_at"`
	NotBefore time.Time `json:"not_before"`
	ExpiresAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (payload *Payload, err error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return
	}
	payload = &Payload{
		Username:  username,
		ID:        tokenId,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
	return
}

var ErrExpiredToken = errors.New("token has expired")

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiresAt) {
		return ErrExpiredToken
	}
	return nil
}
