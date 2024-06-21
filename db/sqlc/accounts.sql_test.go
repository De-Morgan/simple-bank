package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	createParams := CreateAccountParams{
		Owner:    "Victor",
		Balance:  5000,
		Currency: "NGN",
	}
	acct, err := testQueries.CreateAccount(context.Background(), createParams)
	if err != nil {
		log.Fatal("Error occured when creating account", err)
	}
	require.NoError(t, err)
	require.Equal(t, acct.Owner, createParams.Owner)
	require.Equal(t, acct.Balance, createParams.Balance)
	require.Equal(t, acct.Currency, createParams.Currency)
	require.NotZero(t, acct.ID)
	require.NotZero(t, acct.CreatedAt)

}
