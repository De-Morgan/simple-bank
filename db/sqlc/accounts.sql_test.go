package db

import (
	"context"
	"log"
	"testing"

	"github.com/morgan/simplebank/utils"
	"github.com/stretchr/testify/require"
)

var createParams = CreateAccountParams{
	Owner:    utils.RandomName(),
	Balance:  utils.RandomMoney(),
	Currency: "NGN",
}

func createTestAccount() (acct Account) {
	acct, err := testQueries.CreateAccount(context.Background(), createParams)
	if err != nil {
		log.Fatal("Error occured when creating account", err)
	}
	return
}
func TestCreateAccount(t *testing.T) {
	acct := createTestAccount()
	require.Equal(t, acct.Owner, createParams.Owner)
	require.Equal(t, acct.Balance, createParams.Balance)
	require.Equal(t, acct.Currency, createParams.Currency)
	require.NotZero(t, acct.ID)
	require.NotZero(t, acct.CreatedAt)

}
