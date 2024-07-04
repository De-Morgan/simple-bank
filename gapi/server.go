package gapi

import (
	"fmt"

	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/token"
	"github.com/morgan/simplebank/utils"
)

type Server struct {
	pb.UnimplementedSimplebankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config: config,
		store:  store, tokenMaker: tokenMaker}

	return server, nil
}
