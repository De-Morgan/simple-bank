package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: must be at least %d size", minSecretKeySize)
	}
	return PasetoMaker{symmetricKey: []byte(symmetricKey),
		paseto: paseto.NewV2(),
	}, nil
}

func (maker PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)

}

func (maker PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}
	if err = payload.Valid(); err != nil {
		return nil, err
	}
	return &payload, nil

}
