package db

import (
	"context"
	"testing"

	"github.com/morgan/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createUserParams() CreateUserParams {
	hash, _ := utils.HashPassword(utils.RandomString(6))
	return CreateUserParams{
		Username:       utils.RandomName(),
		HashedPassword: hash,
		FullName:       utils.RandomName(),
		Email:          utils.RandomEmail(20),
	}
}
func createTestUser(createUserParams CreateUserParams) (user User, err error) {
	user, err = testQueries.CreateUser(context.Background(), createUserParams)
	return
}
func TestCreateUser(t *testing.T) {
	userParms := createUserParams()
	user, err := createTestUser(userParms)
	require.NoError(t, err)
	require.Equal(t, user.Email, userParms.Email)
	require.Equal(t, user.FullName, userParms.FullName)
	require.NotZero(t, user.Username)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.PasswordChangedAt)

}

func TestGetUserByUsername(t *testing.T) {
	userParms := createUserParams()
	user, err := createTestUser(userParms)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	gUser, err := testQueries.GetUserByUsername(context.Background(), user.Username)
	require.NoError(t, err)
	require.Equal(t, user, gUser)
	require.Equal(t, userParms.Email, gUser.Email)

}
