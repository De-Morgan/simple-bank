// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	AddToAccountBalance(ctx context.Context, arg AddToAccountBalanceParams) (Account, error)
	CreatEntry(ctx context.Context, arg CreatEntryParams) (Entry, error)
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAccount(ctx context.Context, id int64) error
	DeleteEntry(ctx context.Context, id int64) error
	DeleteTransfer(ctx context.Context, id int64) error
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetEntry(ctx context.Context, id int64) (Entry, error)
	GetSession(ctx context.Context, id pgtype.UUID) (Session, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	ListAccount(ctx context.Context, arg ListAccountParams) ([]Account, error)
	ListAccountEntries(ctx context.Context, arg ListAccountEntriesParams) ([]Entry, error)
	ListIncomingTransfers(ctx context.Context, arg ListIncomingTransfersParams) ([]Transfer, error)
	ListOutGoingTransfers(ctx context.Context, arg ListOutGoingTransfersParams) ([]Transfer, error)
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
}

var _ Querier = (*Queries)(nil)
