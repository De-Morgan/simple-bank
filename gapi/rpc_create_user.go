package gapi

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(cxt context.Context, request *pb.CreateUserRequest) (resp *pb.CreateUserResponse, err error) {

	hash, err := utils.HashPassword(request.GetPassword())
	if err != nil {
		err = status.Errorf(codes.Internal, "failed to hash password: %s", err)
		return
	}
	arg := db.CreateUserParams{
		Username:       request.GetUsername(),
		FullName:       request.GetFullName(),
		Email:          request.GetEmail(),
		HashedPassword: hash,
	}
	user, err := server.store.CreateUser(cxt, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "users_email_key", "users_pkey":
				err = status.Errorf(codes.AlreadyExists, "username already exist")
				return
			}
		}
		err = status.Errorf(codes.Internal, "failed to create user %s:%s", arg.Username, err)
		return
	}
	resp = &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return
}
