package access

import (
	"auth/internal/config"
	"auth/internal/service"
)

type serv struct {
	securityConfig config.SecurityConfig
}

func NewService(
	securityConfig config.SecurityConfig,
) service.AccessService {
	return &serv{
		securityConfig: securityConfig,
	}
}
