package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/token"
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
	payload := cxt.MustGet(authorizationPayloadKey).(*token.Payload)

	if account, valid := server.validateAccount(cxt, arg.FromAccountId, request.Currency); !valid {
		return
	} else {
		if account.Owner != payload.Username {
			cxt.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("transfer request declined")))
			return
		}
	}
	if _, valid := server.validateAccount(cxt, arg.ToAccountId, request.Currency); !valid {
		return
	}
	transferSuccess, err := server.store.TransferTx(cxt, arg)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	cxt.JSON(http.StatusCreated, transferSuccess)

}

func (server *Server) validateAccount(cxt *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.getAccount(cxt, accountId)
	if err != nil {
		if err == pgx.ErrNoRows {
			cxt.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}
	if account.Currency == currency {
		return account, true
	}
	cxt.JSON(http.StatusBadRequest, fmt.Errorf("currency mismatch, account currency is %s", account.Currency))
	return account, false
}
