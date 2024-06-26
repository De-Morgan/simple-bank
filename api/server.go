package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/morgan/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccountById)
	router.GET("/accounts", server.ListAccounts)
	router.POST("/transfers", server.TransferMoney)
	router.POST("/users", server.CreateUser)

	server.router = router
	return server
}

func errorResponse(error error) gin.H {
	return gin.H{"error": error.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
