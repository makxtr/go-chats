package access

import (
	"auth/internal/service/access"
	"auth/internal/utils"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, utils.ErrNoMetadata),
		errors.Is(err, utils.ErrNoHeader),
		errors.Is(err, utils.ErrInvalidFormat):
		return status.Error(codes.Unauthenticated, "Authorization token is missing or invalid format")

	case errors.Is(err, access.TokenInvalid):
		return status.Error(codes.Unauthenticated, "Token verification failed (expired or invalid signature)")

	case errors.Is(err, access.Denied):
		return status.Error(codes.PermissionDenied, "Access denied for this role")

	case errors.Is(err, access.RolesGetFailed):
		return status.Error(codes.Internal, "Server failed to retrieve role map")

	default:
		return status.Error(codes.Internal, "Internal server error")
	}
}
