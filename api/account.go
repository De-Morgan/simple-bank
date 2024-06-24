package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/morgan/simplebank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,eq=NGN"`
}

func (server *Server) CreateAccount(cxt *gin.Context) {
	var request CreateAccountRequest
	if err := cxt.ShouldBindBodyWithJSON(&request); err != nil {
		cxt.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    request.Owner,
		Currency: request.Currency,
	}

	account, err := server.store.CreateAccount(cxt, arg)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	cxt.JSON(http.StatusCreated, account)

}

type GetAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) GetAccountById(cxt *gin.Context) {
	var req GetAccountRequest

	if err := cxt.ShouldBindUri(&req); err != nil {
		cxt.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.GetAccount(cxt, req.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			cxt.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("account does not exist")))
			return
		}
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	cxt.JSON(http.StatusOK, account)

}

type ListAccountRequest struct {
	Page  int32
	Limit int32 `binding:"required,min=2,max=50"`
}

func (server *Server) ListAccounts(cxt *gin.Context) {
	var req ListAccountRequest

	if err := cxt.ShouldBindQuery(&req); err != nil {
		cxt.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	args := db.ListAccountParams{
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}
	account, err := server.store.ListAccount(cxt, args)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	cxt.JSON(http.StatusOK, account)

}
