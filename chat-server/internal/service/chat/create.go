package chat

import (
	"chat-server/internal/model"
	"context"

	logModel "github.com/makxtr/go-common/pkg/logger/model"
)

func (s *serv) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.chatRepository.Create(ctx, chat)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &logModel.Log{
			Action:   "chat_created",
			EntityID: id,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
