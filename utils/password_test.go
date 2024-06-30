package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPasswordCorrect(t *testing.T) {
	password := RandomString(6)
	hashPassword, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)
	err = CheckPasswordCorrect(password, hashPassword)
	require.NoError(t, err)

	wrongPassword := RandomString(6)
	err = CheckPasswordCorrect(wrongPassword, hashPassword)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashPassword2, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEqual(t, hashPassword, hashPassword2)

}
