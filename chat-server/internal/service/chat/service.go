package chat

import (
	"chat-server/internal/client/db"
	"chat-server/internal/repository"
	"chat-server/internal/service"
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
