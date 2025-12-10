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

type AuthService interface {
	Login(ctx context.Context, command *model.LoginUserCommand) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, accessToken, endPoint string) (bool, error)
}
