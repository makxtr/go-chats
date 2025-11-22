package user

import (
	"auth/internal/repository"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, repository.ErrNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, repository.ErrQueryBuild):
		return status.Error(codes.Internal, "internal error: failed to build query")
	case errors.Is(err, repository.ErrCreateFailed):
		return status.Error(codes.Internal, "failed to create user")
	case errors.Is(err, repository.ErrUpdateFailed):
		return status.Error(codes.Internal, "failed to update user")
	case errors.Is(err, repository.ErrDeleteFailed):
		return status.Error(codes.Internal, "failed to delete user")
	default:
		return status.Error(codes.Internal, err.Error())
	}
}