package user

import (
	"auth/internal/model"
	"context"
)

func (s *serv) Update(ctx context.Context, id int64, updateUser *model.UpdateUserData) error {
	err := s.userRepository.Update(ctx, id, updateUser)
	if err != nil {
		return err
	}

	return nil
}
