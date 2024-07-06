package gapi

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/morgan/simplebank/db/sqlc"
	"github.com/morgan/simplebank/pb"
	"github.com/morgan/simplebank/utils"
	"github.com/morgan/simplebank/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(cxt context.Context, request *pb.UpdateUserRequest) (resp *pb.UpdateUserResponse, err error) {
	if violations := validateUpdateUserRequest(request); violations != nil {
		return nil, invalidArguementError(violations)
	}
	payload, err := server.authorizeUser(cxt)
	if err != nil {
		return nil, unAuthenticatedError(err)
	}
	if payload.Username != request.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "can not update other user's info")
	}
	password := request.GetData().GetPassword()
	var hashedPass string
	if password != "" {
		hashedPass, err = utils.HashPassword(password)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "unable to hash passsword")
		}
	}

	arg := db.UpdateUserParams{
		FullName: pgtype.Text{
			String: request.GetData().GetFullName(),
			Valid:  request.GetData().GetFullName() != "",
		},
		HashedPassword: pgtype.Text{
			String: hashedPass,
			Valid:  hashedPass != "",
		},
		Email: pgtype.Text{
			String: request.GetData().GetEmail(),
			Valid:  request.GetData().GetEmail() != "",
		},
		Username: request.GetUsername(),
	}
	user, err := server.store.UpdateUser(cxt, arg)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.ConstraintName {
			case "users_email_key":
				err = status.Errorf(codes.AlreadyExists, "email already in use")
				return
			}
		}
		if err == pgx.ErrNoRows {
			err = status.Errorf(codes.NotFound, "user not found")
			return
		}
		err = status.Errorf(codes.Internal, "failed to Update user %s:%s", arg.Username, err)
		return
	}
	resp = &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return
}

func validateUpdateUserRequest(request *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidEmailAddress(request.GetData().GetEmail()); err != nil && request.GetData().GetEmail() != "" {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := validation.ValidatePassword(request.GetData().GetPassword()); err != nil && request.GetData().GetPassword() != "" {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := validation.ValidateFullName(request.GetData().GetFullName()); err != nil && request.GetData().GetFullName() != "" {
		violations = append(violations, fieldViolation("full_name", err))
	}
	return
}
