package token

// import (
// 	"fmt"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// const minSecretKeySize = 32

// // JWTMaker is a JSON Web Token maker
// type JWTMaker struct {
// 	secretKey string
// }

// func (jwtMaker JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
// 	payload, err := NewPayload(username, duration)
// 	if err != nil {
// 		return "", err
// 	}

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
// 	return jwtToken.SignedString([]byte(jwtMaker.secretKey))
// }

// func (jwtMaker JWTMaker) VerifyToken(token string) (*Payload, error) {
// 	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.ErrTokenSignatureInvalid
// 		}
// 		return []byte(jwtMaker.secretKey), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	payload, ok := jwtToken.Claims.(*Payload)
// 	if !ok {
// 		return nil, jwt.ErrTokenInvalidClaims
// 	}

// 	return payload, nil

// }

// func NewJWTMaker(secretKey string) (Maker, error) {
// 	if len(secretKey) < minSecretKeySize {
// 		return nil, fmt.Errorf("invalid secret key size: must be at least %d size", minSecretKeySize)
// 	}
// 	return JWTMaker{secretKey: secretKey}, nil
// }
