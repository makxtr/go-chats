package chat

import (
	"chat-server/internal/model"
	"context"

	"github.com/pkg/errors"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		errTx = s.chatRepository.Delete(ctx, id)
		if errTx != nil {
			return errors.Wrap(errTx, "failed to delete chat from repository")
		}

		errTx = s.logRepository.Log(ctx, &model.ChatLog{
			Action:   "chat_deleted",
			EntityID: id,
		})
		if errTx != nil {
			return errors.Wrap(errTx, "failed to log chat deletion")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "failed to delete chat with transaction")
	}

	return nil
}
