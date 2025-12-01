package user_test

import (
	"context"
	"errors"
	"testing"

	"auth/internal/api/user"
	"auth/internal/repository/mocks"
	userService "auth/internal/service/user"
	desc "auth/pkg/user_v1"

	"github.com/gojuno/minimock/v3"
	logModel "github.com/makxtr/go-common/pkg/logger/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestImplementation_Delete(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) *mocks.UserRepositoryMock
	type logRepositoryMockFunc func(mc *minimock.Controller) *mocks.LogRepositoryMock

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id = int64(1)

		req = &desc.DeleteRequest{
			Id: id,
		}

		logEntry = &logModel.Log{
			Action:   "user_deleted",
			EntityID: id,
		}

		repoErr = errors.New("repository error")
		logErr  = errors.New("log error")
	)

	tests := []struct {
		name               string
		args               args
		want               *emptypb.Empty
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
			want: &emptypb.Empty{},
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
			want: nil,
			err:  repoErr,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(repoErr)
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
			want: nil,
			err:  logErr,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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

			resp, err := api.Delete(tt.args.ctx, tt.args.req)

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
