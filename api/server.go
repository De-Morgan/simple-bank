package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/morgan/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.CreateAccount)
	router.GET("/accounts/:id", server.GetAccountById)
	router.GET("/accounts", server.ListAccounts)

	server.router = router
	return server
}

func errorResponse(error error) gin.H {
	return gin.H{"error": error.Error()}
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
