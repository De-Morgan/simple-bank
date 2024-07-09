package gapi

import (
	"fmt"

	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/token"
	"github.com/morgan/simplebank/utils"
	"github.com/morgan/simplebank/worker"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config          utils.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

func NewServer(config utils.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		taskDistributor: taskDistributor,
		store:           store, tokenMaker: tokenMaker}

	return server, nil
}
