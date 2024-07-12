package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
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

func TestUpdateUser(t *testing.T) {
	nName := utils.RandomName()
	nEmail := utils.RandomEmail(6)
	nHashpassword, err := utils.HashPassword(utils.RandomString(6))
	require.NoError(t, err)
	userParms := createUserParams()
	user, err := createTestUser(userParms)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	tests := []struct {
		name     string
		arg      UpdateUserParams
		validate func(t *testing.T, newUser User, err error)
	}{

		{
			name: "Update Fullname",
			arg: UpdateUserParams{
				FullName: pgtype.Text{
					String: nName,
					Valid:  true,
				},
				Username: user.Username,
			},
			validate: func(t *testing.T, newUser User, err error) {
				require.NoError(t, err)
				require.Equal(t, user.Email, newUser.Email)
				require.Equal(t, user.HashedPassword, newUser.HashedPassword)
				require.NotEqual(t, user.FullName, newUser.FullName)
				require.Equal(t, newUser.FullName, nName)
			},
		},
		{
			name: "Update Email",
			arg: UpdateUserParams{
				Email: pgtype.Text{
					String: nEmail,
					Valid:  true,
				},
				Username: user.Username,
			},
			validate: func(t *testing.T, newUser User, err error) {
				require.NoError(t, err)
				require.Equal(t, newUser.Email, nEmail)
			},
		},
		{
			name: "Update All Fields",
			arg: UpdateUserParams{
				Email: pgtype.Text{
					String: nEmail,
					Valid:  true,
				},
				FullName: pgtype.Text{
					String: nName,
					Valid:  true,
				},
				HashedPassword: pgtype.Text{
					String: nHashpassword,
					Valid:  true,
				},
				Username: user.Username,
			},
			validate: func(t *testing.T, newUser User, err error) {
				require.NoError(t, err)
				require.Equal(t, newUser.Email, nEmail)
				require.Equal(t, newUser.FullName, nName)
				require.Equal(t, newUser.HashedPassword, nHashpassword)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			uUser, err := testQueries.UpdateUser(context.Background(), tt.arg)

			tt.validate(t, uUser, err)

		})
	}

}
