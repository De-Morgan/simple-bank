package db

import (
	"context"
	"log"
	"testing"

	"github.com/morgan/simplebank/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createParams() CreateAccountParams {
	userCreateParms := createUserParams()
	user, _ := createTestUser(userCreateParms)
	return CreateAccountParams{
		Owner:    user.Username,
		Balance:  1000,
		Currency: utils.RandomCurrency(),
	}
}

func createTestAccount() (acct Account) {
	acct, err := testQueries.CreateAccount(context.Background(), createParams())
	if err != nil {
		log.Fatal("Error occured when creating account", err)
	}
	return
}
func TestCreateAccount(t *testing.T) {
	acct := createTestAccount()
	require.NotZero(t, acct.ID)
	require.NotZero(t, acct.CreatedAt)

}

func TestDeleteAccount(t *testing.T) {
	acct := createTestAccount()
	require.NotEmpty(t, acct)
	err := testQueries.DeleteAccount(context.Background(), acct.ID)
	require.NoError(t, err)
	acct2, err := testQueries.GetAccount(context.Background(), acct.ID)
	require.Error(t, err)
	require.Empty(t, acct2)
}

func TestGetAccount(t *testing.T) {
	acct := createTestAccount()
	require.NotEmpty(t, acct)
	gAcct, err := testQueries.GetAccount(context.Background(), acct.ID)
	require.NoError(t, err)
	assert.Equal(t, acct, gAcct)
}

func TestListAccount(t *testing.T) {
	acct := createTestAccount()
	require.NotEmpty(t, acct)
	accts, err := testQueries.ListAccount(context.Background(), ListAccountParams{1, 1})
	require.NoError(t, err)
	assert.NotEmpty(t, accts)
	require.Len(t, accts, 1)

}

func TestUpdateAccount(t *testing.T) {
	acct := createTestAccount()
	args := UpdateAccountParams{
		ID: acct.ID, Balance: utils.RandomMoney(),
	}
	acct, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	assert.NotEmpty(t, acct)
	assert.Equal(t, acct.Balance, args.Balance)
}
