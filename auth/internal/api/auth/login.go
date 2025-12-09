package auth

import (
	"auth/internal/converter"
	desc "auth/pkg/auth_v1"
	"context"
	"log"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	token, err := i.authService.Login(ctx, converter.ToLoginFromDesc(req))
	if err != nil {
		return nil, mapError(err)
	}

	log.Printf("login user with username: %s", req.Username)

	return &desc.LoginResponse{RefreshToken: token.Token}, nil
}
