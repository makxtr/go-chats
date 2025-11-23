package repository

import (
	"auth/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, createUser *model.CreateUserData) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, updateUser *model.UpdateUserData) error
}

type LogRepository interface {
	Log(ctx context.Context, userLog *model.UserLog) error
}
