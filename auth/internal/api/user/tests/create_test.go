package user_test

import (
	"context"
	"errors"
	"testing"

	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/repository/mocks"
	userService "auth/internal/service/user"
	desc "auth/pkg/user_v1"

	"github.com/gojuno/minimock/v3"
	logModel "github.com/makxtr/go-common/pkg/logger/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImplementation_Create(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) *mocks.UserRepositoryMock
	type logRepositoryMockFunc func(mc *minimock.Controller) *mocks.LogRepositoryMock

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id              = int64(1)
		name            = "test_user"
		email           = "test@example.com"
		password        = "password123"
		passwordConfirm = "password123"
		role            = desc.Role_ROLE_USER

		req = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}

		// Note: We can't easily predict the exact CreateUserCommand because of password hashing
		// So we will use minimock.Any for the command argument in the mock expectation
		// or we could match fields if we really wanted to, but Any is safer for now given hashing.

		logEntry = &logModel.Log{
			Action:   "user_created",
			EntityID: id,
		}

		repoErr = errors.New("repository error")
		logErr  = errors.New("log error")
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Set(func(ctx context.Context, createUser *model.CreateUserData) (int64, error) {
					assert.Equal(t, name, createUser.Info.Name)
					assert.Equal(t, email, createUser.Info.Email)
					assert.Equal(t, model.RoleUser, createUser.Info.Role)
					assert.NotEmpty(t, createUser.HashedPassword)
					return id, nil
				})
				mock.GetMock.Expect(ctx, id).Return(nil, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				mock := mocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(ctx, logEntry).Return(nil)
				return mock
			},
		},
		{
			name: "repository error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Set(func(ctx context.Context, createUser *model.CreateUserData) (int64, error) {
					return 0, repoErr
				})
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				mock := mocks.NewLogRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "log error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: 0,
			err:  logErr,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Set(func(ctx context.Context, createUser *model.CreateUserData) (int64, error) {
					return id, nil
				})
				mock.GetMock.Expect(ctx, id).Return(nil, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				mock := mocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(ctx, logEntry).Return(logErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := tt.userRepositoryMock(mc)
			logRepoMock := tt.logRepositoryMock(mc)

			txManager := &txManagerMock{}

			service := userService.NewService(
				userRepoMock,
				logRepoMock,
				txManager,
			)

			api := user.NewImplementation(service)

			resp, err := api.Create(tt.args.ctx, tt.args.req)

			if tt.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.want, resp.Id)
			}
		})
	}
}
