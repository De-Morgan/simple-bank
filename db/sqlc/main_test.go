package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	if conn, err := pgx.Connect(context.Background(), dbSource); err != nil {
		log.Fatal("Can't connect to db", err)
	} else {
		testQueries = New(conn)
	}
	os.Exit(m.Run())
}
