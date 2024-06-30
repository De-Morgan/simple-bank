package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morgan/simplebank/api"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatalln("Can't load config files", err)
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Can't start server")
	}
}
