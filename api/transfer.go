package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/morgan/simplebank/db/sqlc"
)

type TransferRequest struct {
	FromAccountID int64  `json:"fromAccountId" binding:"required,min=1"`
	ToAccountID   int64  `json:"toAccountId" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) TransferMoney(cxt *gin.Context) {
	var request TransferRequest
	if err := cxt.ShouldBindBodyWithJSON(&request); err != nil {
		cxt.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.TransferParam{
		FromAccountId: request.FromAccountID,
		ToAccountId:   request.ToAccountID,
		Amount:        request.Amount,
	}
	if !server.checkAccountCurrency(cxt, request.Currency, arg.FromAccountId) {
		return
	}
	if !server.checkAccountCurrency(cxt, request.Currency, arg.ToAccountId) {
		return
	}
	transferSuccess, err := server.store.TransferTx(cxt, arg)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cxt.JSON(http.StatusCreated, transferSuccess)

}

// Returns error if currency doesn't match account currency
func (server *Server) checkAccountCurrency(cxt *gin.Context, currency string, accountId int64) bool {
	account, err := server.getAccount(cxt, accountId)
	if err != nil {
		return false
	}
	if account.Currency == currency {
		return true
	}
	cxt.JSON(http.StatusBadRequest, fmt.Errorf("currency mismatch, account currency is %s", account.Currency))
	return false
}
