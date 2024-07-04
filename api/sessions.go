package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/morgan/simplebank/db/sqlc"
)

type CreateSessionRequest struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	Username     string    `json:"username" binding:"required"`
	RefreshToken string    `json:"refresh_token" binding:"required"`
	IsBlocked    bool      `json:"is_blocked" binding:"required"`
	ExpiresAt    time.Time `json:"expires_at" binding:"required"`
}

func (server *Server) createSession(cxt *gin.Context, request CreateSessionRequest) (session db.Session, err error) {

	arg := db.CreateSessionParams{
		ID:           pgtype.UUID{Bytes: request.ID, Valid: true},
		Username:     request.Username,
		RefreshToken: request.RefreshToken,
		UserAgent:    cxt.Request.UserAgent(),
		ClientIp:     cxt.ClientIP(),
		IsBlocked:    request.IsBlocked,
		ExpiresAt:    pgtype.Timestamptz{Time: request.ExpiresAt, Valid: true},
	}
	session, err = server.store.CreateSession(cxt, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "sessions_username_fkey":
				err = fmt.Errorf("user not found")
				return
			}
		}
		cxt.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	return
}
