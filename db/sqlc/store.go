package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// / Store provides all functions to execute individual query or transaction
type Store struct {
	*Queries
	conn *pgxpool.Pool
}

func NewStore(conn *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(conn),
		conn:    conn,
	}
}

func (store *Store) execTx(cxt context.Context, fn func(*Queries) error) error {
	tx, err := store.conn.Begin(cxt)
	if err != nil {
		return err
	}
	query := New(tx)
	err = fn(query)
	if err != nil {
		if e := tx.Rollback(cxt); e != nil {
			return fmt.Errorf("tx Err %s, db error %s", e, err)
		}
		return err
	}
	return tx.Commit(cxt)
}
