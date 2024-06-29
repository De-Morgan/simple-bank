package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/token"
	"github.com/morgan/simplebank/utils"
)

type Server struct {
	config     utils.Config
	store      db.Store
	router     *gin.Engine
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

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}
	server.setUpRounter()
	return server, nil
}

func (server *Server) setUpRounter() {
	router := gin.Default()
	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccountById)
	router.GET("/accounts", server.ListAccounts)
	router.POST("/transfers", server.TransferMoney)
	router.POST("/users", server.CreateUser)
	router.POST("/users/login", server.LoginUser)
	server.router = router
}

func errorResponse(error error) gin.H {
	return gin.H{"error": error.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
