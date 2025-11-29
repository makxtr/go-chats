package user

import (
	"auth/internal/repository"
	"auth/internal/service"

	"github.com/makxtr/go-common/pkg/db"
)

type serv struct {
	userRepository repository.UserRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

func NewService(
	userRepository repository.UserRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.UserService {
	return &serv{
		userRepository: userRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
