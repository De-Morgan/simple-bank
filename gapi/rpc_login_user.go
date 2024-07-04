package gapi

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(cxt context.Context, request *pb.LoginUserRequest) (resp *pb.LoginUserResponse, err error) {
	password := request.Password
	user, err := server.store.GetUserByUsername(cxt, request.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = status.Errorf(codes.NotFound, "User not found")
			return
		}
		err = status.Errorf(codes.Internal, "Error getting user: %s", err)
		return
	}
	userHashedPassword := user.HashedPassword
	if err = utils.CheckPasswordCorrect(password, userHashedPassword); err != nil {
		err = status.Errorf(codes.NotFound, "invalid credential ")
		return
	}
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		err = status.Errorf(codes.Internal, "Error creating accessToken: %s", err)
		return
	}
	refreshToken, payload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration)
	if err != nil {
		err = status.Errorf(codes.Internal, "Error creating refeshToken: %s", err)
		return
	}
	_, err = server.store.CreateSession(cxt, db.CreateSessionParams{
		ID: pgtype.UUID{
			Bytes: payload.ID,
			Valid: true,
		},
		Username:     user.Username,
		RefreshToken: refreshToken,
		ExpiresAt: pgtype.Timestamptz{
			Time: payload.ExpiresAt, Valid: true,
		},
		IsBlocked: false,
	})
	if err != nil {
		err = status.Errorf(codes.Internal, "Error creating session: %s", err)
		return
	}

	resp = &pb.LoginUserResponse{
		SessionId:             payload.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(payload.ExpiresAt),
		User:                  convertUser(user),
	}
	return
}
