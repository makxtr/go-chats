package auth

import (
	"auth/internal/service/auth"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, auth.TokenGenFailed):
		return status.Error(codes.Internal, "invalid credentials provided")
	case errors.Is(err, auth.LoginFailed):
		return status.Error(codes.Unauthenticated, "failed to login user")
	default:
		return status.Error(codes.Internal, err.Error())
	}

}
