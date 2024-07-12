package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type refreshTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) refreshToken(cxt *gin.Context) {
	var request refreshTokenRequest
	if err := cxt.ShouldBindBodyWithJSON(&request); err != nil {
		cxt.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	refreshPayload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		cxt.JSON(http.StatusUnauthorized, errorResponse(err))
	}
	session, err := server.store.GetSession(cxt, pgtype.UUID{
		Bytes: refreshPayload.ID,
		Valid: true,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			cxt.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if session.IsBlocked {
		cxt.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("blocked session")))
		return
	}
	if session.Username != refreshPayload.Username {
		cxt.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("incorrect session user")))
		return
	}

	if session.RefreshToken != request.RefreshToken {
		cxt.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("mismatch session token")))
		return
	}

	if time.Now().After(session.ExpiresAt.Time) {
		cxt.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("expired session")))
		return
	}

	nAccessToken, nPayload, err := server.tokenMaker.CreateToken(session.Username, server.config.AccessTokenDuration)
	if err != nil {
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := refreshTokenResponse{
		AccessToken:          nAccessToken,
		AccessTokenExpiresAt: nPayload.ExpiresAt,
	}
	cxt.JSON(http.StatusOK, res)
}
