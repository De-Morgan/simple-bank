package db

import (
	"context"
)

type moneyTransferKey string

var AmountSentKey moneyTransferKey = "AmountSent"

type TransferParam struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}
type TransferResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(cxt context.Context, arg TransferParam) (TransferResult, error) {
	var result TransferResult

	err := store.execTx(cxt, func(q *Queries) error {

		var err error
		result.Transfer, err = q.CreateTransfer(cxt, CreateTransferParams{FromAccountID: arg.FromAccountId, ToAccountID: arg.ToAccountId, Amount: arg.Amount})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreatEntry(cxt, CreatEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreatEntry(cxt, CreatEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromAccount, result.ToAccount, err = addMoney(cxt, q, arg.FromAccountId, arg.ToAccountId, arg.Amount)
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}

func addMoney(cxt context.Context, q *Queries, fromAccountId, toAccountId, money int64) (fromAccount, toAccount Account, err error) {
	//Update account with smaller id first to avoid deadlock
	if fromAccountId < toAccountId {
		fromAccount, err = q.AddToAccountBalance(cxt, AddToAccountBalanceParams{
			ID:     fromAccountId,
			Amount: -money,
		})
		if err != nil {
			return
		}
		toAccount, err = q.AddToAccountBalance(cxt, AddToAccountBalanceParams{ID: toAccountId, Amount: money})
		if err != nil {
			return
		}
	} else {
		toAccount, err = q.AddToAccountBalance(cxt, AddToAccountBalanceParams{ID: toAccountId, Amount: money})
		if err != nil {
			return
		}
		fromAccount, err = q.AddToAccountBalance(cxt, AddToAccountBalanceParams{
			ID:     fromAccountId,
			Amount: -money,
		})
		if err != nil {
			return
		}
	}

	return
}
