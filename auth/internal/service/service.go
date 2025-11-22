package service

import (
	"auth/internal/model"
	"context"
)

type UserService interface {
	Create(ctx context.Context, command *model.CreateUserCommand) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, updateUser *model.UpdateUserData) error
	Delete(ctx context.Context, id int64) error
}
