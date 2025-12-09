package auth

import (
	"auth/internal/config"
	"auth/internal/repository"
	"auth/internal/service"

	"github.com/makxtr/go-common/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
	securityConfig config.SecurityConfig
}

func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
	securityConfig config.SecurityConfig,
) service.AuthService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
		securityConfig: securityConfig,
	}
}
