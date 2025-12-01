package user_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"auth/internal/api/user"
	"auth/internal/model"
	"auth/internal/repository/mocks"
	userService "auth/internal/service/user"
	desc "auth/pkg/user_v1"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestImplementation_Get(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) *mocks.UserRepositoryMock
	type logRepositoryMockFunc func(mc *minimock.Controller) *mocks.LogRepositoryMock

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = int64(1)
		name      = "test_user"
		email     = "test@example.com"
		role      = desc.Role_ROLE_USER
		createdAt = time.Now()
		updatedAt = time.Now()

		req = &desc.GetRequest{
			Id: id,
		}

		userModel = &model.User{
			ID: id,
			Info: model.UserInfo{
				Name:  name,
				Email: email,
				Role:  model.RoleUser,
			},
			CreatedAt: createdAt,
			UpdatedAt: sql.NullTime{
				Time:  updatedAt,
				Valid: true,
			},
		}

		res = &desc.GetResponse{
			User: &desc.User{
				Id:        id,
				Name:      name,
				Email:     email,
				Role:      role,
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}

		repoErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		want               *desc.GetResponse
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
			want: res,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(userModel, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				mock := mocks.NewLogRepositoryMock(mc)
				return mock
			},
		},
		{
			name: "repository error",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				mock := mocks.NewLogRepositoryMock(mc)
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

			resp, err := api.Get(tt.args.ctx, tt.args.req)

			if tt.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tt.want, resp)
			}
		})
	}
}
