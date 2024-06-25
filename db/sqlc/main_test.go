package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morgan/simplebank/utils"
)

var testQueries *Queries
var testDbConn *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Printf("Can't log config file: %v\n", err)
		os.Exit(1)
	}

	testDbConn, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer testDbConn.Close()
	testQueries = New(testDbConn)
	os.Exit(m.Run())
}
