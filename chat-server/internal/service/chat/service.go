package chat

import (
	"chat-server/internal/repository"
	"chat-server/internal/service"

	"github.com/makxtr/go-common/pkg/db"
)

type serv struct {
	chatRepository repository.ChatRepository
	logRepository  repository.LogRepository
	txManager      db.TxManager
}

func NewService(
	chatRepository repository.ChatRepository,
	logRepository repository.LogRepository,
	txManager db.TxManager,
) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		logRepository:  logRepository,
		txManager:      txManager,
	}
}
