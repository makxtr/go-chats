package auth

import (
	"auth/internal/model"
	"auth/internal/utils"
	"context"
	"time"
)

func (s *serv) Login(ctx context.Context, command *model.LoginUserCommand) (string, error) {
	userSecure, err := s.userRepository.GetForLogin(ctx, command.Username)
	if err != nil {
		return "", err
	}
	if !utils.VerifyPassword(userSecure.Password, command.Password) {
		return "", LoginFailed
	}

	refreshToken, err := utils.GenerateToken(model.UserInfo{
		Name: command.Username,
		Role: userSecure.Role,
	},
		[]byte(s.securityConfig.RefreshKey()),
		time.Duration(s.securityConfig.RefreshExp())*time.Minute,
	)
	if err != nil {
		return "", TokenGenFailed
	}

	return refreshToken, nil
}
