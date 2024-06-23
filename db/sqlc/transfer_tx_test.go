package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDbConn)

	acct1 := createTestAccount()
	acct2 := createTestAccount()

	money := []int64{50, 100, 40, 10}
	results := make(chan TransferResult, len(money))

	for _, m := range money {
		go func(amount int64) {
			r, _ := store.TransferTx(context.Background(), TransferParam{
				FromAccountId: acct1.ID,
				ToAccountId:   acct2.ID,
				Amount:        amount,
			})
			results <- r

		}(m)
	}

	for range money {
		result := <-results
		require.NotEmpty(t, result)
		require.NotZero(t, result)

		//Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		require.Equal(t, result.Transfer.FromAccountID, acct1.ID)
		require.Equal(t, result.Transfer.ToAccountID, acct2.ID)
		_, err := store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Check entries
		fromEntry := result.FromEntry
		toEntry := result.ToEntry

		require.NotEmpty(t, fromEntry)
		require.NotEmpty(t, toEntry)
		require.NotZero(t, fromEntry.CreatedAt)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, acct1.ID)
		require.Equal(t, toEntry.AccountID, acct2.ID)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//Check account and balance
		FromAccount := result.FromAccount
		ToAccount := result.ToAccount
		require.NotEmpty(t, FromAccount)
		require.NotEmpty(t, ToAccount)
		require.Equal(t, FromAccount.ID, acct1.ID)
		require.Equal(t, ToAccount.ID, acct2.ID)

	}

}

func TestForDeadLockTransferTx(t *testing.T) {
	store := NewStore(testDbConn)
	acct1 := createTestAccount()
	acct2 := createTestAccount()
	money := []int64{10, 10, 10, 10, 10, 10}
	errors := make(chan error, len(money))

	for i, m := range money {
		fromId := acct1.ID
		toId := acct2.ID
		if i%2 == 0 {
			fromId = acct2.ID
			toId = acct1.ID
		}
		go func(amount int64) {
			_, err := store.TransferTx(context.Background(), TransferParam{
				FromAccountId: fromId,
				ToAccountId:   toId,
				Amount:        amount,
			})
			errors <- err
		}(m)
	}

	for range money {
		err := <-errors
		require.NoError(t, err)
	}

	updatedAcct1, _ := store.GetAccount(context.Background(), acct1.ID)
	require.Equal(t, acct1.Balance, updatedAcct1.Balance)
	updatedAcct2, _ := store.GetAccount(context.Background(), acct2.ID)
	require.Equal(t, acct2.Balance, updatedAcct2.Balance)

}
