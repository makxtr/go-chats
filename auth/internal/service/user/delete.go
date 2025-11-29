package user

import (
	"context"

	logModel "github.com/makxtr/go-common/pkg/logger/model"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Delete(ctx, id)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &logModel.Log{
			Action:   "user_deleted",
			EntityID: id,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
