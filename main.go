package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/morgan/simplebank/api"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/gapi"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	//runGrpcServer(config, store)
	runGatewayServer(config, store)
}

func runGrpcServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't start server")
	}
	grpcserver := grpc.NewServer()

	pb.RegisterSimplebankServer(grpcserver, server)
	reflection.Register(grpcserver)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalln("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcserver.Serve(listener)
	if err != nil {
		log.Fatalln("cannot start gRPC server")
	}

}

func runGatewayServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Can't start server")
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
		log.Fatalln("cannot start gRPC server: ", err)
	}
	fs := http.FileServer(http.Dir("./doc/swagger"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatalln("cannot register handler server: ", err)
	}
	log.Printf("start Http gateway server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalln("cannot start http gateway server")
	}
}

func runGinServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	err = server.Start(config.HttpServerAddress)

	if err != nil {
		log.Fatal("Can't start server")
	}
}
