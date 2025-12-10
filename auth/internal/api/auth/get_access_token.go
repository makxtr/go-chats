package auth

import (
	"context"
	"log"

	desc "auth/pkg/auth_v1"
)

func (i *Implementation) GetAccessToken(ctx context.Context, req *desc.GetAccessTokenRequest) (*desc.GetAccessTokenResponse, error) {
	log.Printf("getting access token")

	accessToken, err := i.authService.GetAccessToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, mapError(err)
	}

	return &desc.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}