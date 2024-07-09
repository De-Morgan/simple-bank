package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morgan/simplebank/api"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/gapi"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/utils"
	"github.com/morgan/simplebank/worker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Can't load config files")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Run db migration\
	runDBMigration(config.MigrationUrl, config.DBSource)

	store := db.NewStore(conn)
	redisClientOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisDistributor(redisClientOpt)
	go runTaskProcessor(redisClientOpt, store)
	go runGrpcServer(config, store, taskDistributor)
	runGatewayServer(config, store, taskDistributor)
}

func runDBMigration(migrationUrl, dbSource string) {
	m, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot create a new migration instance: ")
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().
			Err(err).
			Msg("failed to run migrate up: ")
	}
	log.Info().Msg("db migrated successfully")
}
func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(&redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start processor")
	}
}

func runGrpcServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Can't start server")
	}
	option := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcserver := grpc.NewServer(option)

	pb.RegisterSimplebankServer(grpcserver, server)
	reflection.Register(grpcserver)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().
			Msg("cannot create listener")
	}
	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

}

func runGatewayServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("Can't start server")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpc := runtime.NewServeMux(jsonOption)
	pb.RegisterSimplebankHandlerServer(ctx, grpc, server)

	mux := http.NewServeMux()
	mux.Handle("/", grpc)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot start gRPC server: ")
	}
	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("cannot register handler server: ")
	}
	log.Info().Msgf("start Http gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start http gateway server")
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal().Msg("Can't start server")
	}
}
