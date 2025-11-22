package service

import (
	"chat-server/internal/model"
	"context"
)

type ChatService interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, message *model.Message) error
}
