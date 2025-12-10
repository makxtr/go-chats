package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"time"
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.securityConfig.RefreshKey()))
	if err != nil {
		return "", RefreshTokenInvalid
	}

	role := model.RoleUnspecified
	switch claims.Role {
	case "USER":
		role = model.RoleUser
	case "ADMIN":
		role = model.RoleAdmin
	}

	accessToken, err := utils.GenerateToken(
		model.UserInfo{
			Name: claims.Username,
			Role: role,
		},
		[]byte(s.securityConfig.AccessKey()),
		time.Duration(s.securityConfig.AccessExp())*time.Minute,
	)
	if err != nil {
		return "", TokenGenFailed
	}

	return accessToken, nil
}
