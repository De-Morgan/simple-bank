package token

import (
	"testing"
	"time"

	"github.com/morgan/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoCreateToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)
	username := utils.RandomName()
	duration := time.Minute

	//issueAt := time.Now()
	//expiredAt := issueAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)

	require.Equal(t, username, payload.Username)

}

func TestPasetoExpiredJWTToken(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)
	username := utils.RandomName()
	duration := time.Second

	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	_, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.Equal(t, err.Error(), ErrExpiredToken.Error())

}
