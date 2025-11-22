package converter

import (
	"chat-server/internal/model"
	modelRepo "chat-server/internal/repository/chat/model"
)

func ToChatFromRepo(chat *modelRepo.Chat) *model.Chat {
	return &model.Chat{
		ID:        chat.ID,
		Usernames: chat.Usernames,
	}
}
