package user

import (
	"auth/internal/model"
	"context"

	logModel "github.com/makxtr/go-common/pkg/logger/model"
)

func (s *serv) Update(ctx context.Context, id int64, updateUser *model.UpdateUserData) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		errTx := s.userRepository.Update(ctx, id, updateUser)
		if errTx != nil {
			return errTx
		}

		errTx = s.logRepository.Log(ctx, &logModel.Log{
			Action:   "user_updated",
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
