package token

// import (
// 	"testing"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// 	"github.com/morgan/simplebank/utils"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreateToken(t *testing.T) {
// 	maker, err := NewJWTMaker(utils.RandomString(32))
// 	require.NoError(t, err)
// 	username := utils.RandomName()
// 	duration := time.Minute

// 	//issueAt := time.Now()
// 	//expiredAt := issueAt.Add(duration)

// 	token, err := maker.CreateToken(username, duration)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	payload, err := maker.VerifyToken(token)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, payload)
// 	require.NotZero(t, payload.ID)

// 	require.Equal(t, username, payload.Username)

// }

// func TestExpiredJWTToken(t *testing.T) {
// 	maker, err := NewJWTMaker(utils.RandomString(32))
// 	require.NoError(t, err)
// 	username := utils.RandomName()
// 	duration := time.Second

// 	token, err := maker.CreateToken(username, -duration)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)
// 	_, err = maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.Contains(t, err.Error(), jwt.ErrTokenExpired.Error())

// }

// func TestInvalidJwtAlgNonToken(t *testing.T) {
// 	payload, err := NewPayload(utils.RandomName(), time.Minute)
// 	require.NoError(t, err)

// 	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
// 	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, token)

// 	maker, err := NewJWTMaker(utils.RandomString(32))
// 	require.NoError(t, err)
// 	payload, err = maker.VerifyToken(token)
// 	require.Error(t, err)
// 	require.Nil(t, payload)
// 	require.Contains(t, err.Error(), jwt.ErrTokenSignatureInvalid.Error())

// }
