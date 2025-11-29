package repository

import (
	"chat-server/internal/model"
	"context"

	logModel "github.com/makxtr/go-common/pkg/logger/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type LogRepository interface {
	Log(ctx context.Context, log *logModel.Log) error
}
