package gapi

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/utils"
	"github.com/morgan/simplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(cxt context.Context, request *pb.CreateUserRequest) (resp *pb.CreateUserResponse, err error) {
	if violations := validateCreateUserRequest(request); violations != nil {
		return nil, invalidArguementError(violations)
	}
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
			case "users_email_key":
				err = status.Errorf(codes.AlreadyExists, "email already in use")
				return
			case "users_pkey":
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

func validateCreateUserRequest(request *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidatePassword(request.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := validation.ValidateFullName(request.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := validation.ValidEmailAddress(request.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return
}