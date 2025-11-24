package repository

import (
	"chat-server/internal/model"
	"context"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type LogRepository interface {
	Log(ctx context.Context, chatLog *model.ChatLog) error
}
