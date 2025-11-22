package chat

import (
	"chat-server/internal/service"
	desc "chat-server/pkg/chat_server_v1"
)

type Implementation struct {
	desc.UnimplementedChatServerV1Server
	chatService service.ChatService
}

func NewImplementation(chatService service.ChatService) *Implementation {
	return &Implementation{
		chatService: chatService,
	}
}
