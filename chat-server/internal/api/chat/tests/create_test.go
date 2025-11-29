package chat_test

import (
	"chat-server/internal/api/chat"
	"chat-server/internal/model"
	"chat-server/internal/repository/mocks"
	chatService "chat-server/internal/service/chat"
	desc "chat-server/pkg/chat_server_v1"
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	logModel "github.com/makxtr/go-common/pkg/logger/model"
	"github.com/stretchr/testify/require"
)

func TestImplementation_Create(t *testing.T) {
	type chatRepositoryMockFunc func(mc *minimock.Controller) *mocks.ChatRepositoryMock
	type logRepositoryMockFunc func(mc *minimock.Controller) *mocks.LogRepositoryMock

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		usernames = []string{"user1", "user2"}
		chatID    = int64(123)

		req = &desc.CreateRequest{
			Usernames: usernames,
		}

		chatModel = &model.Chat{
			Usernames: usernames,
		}

		logEntry = &logModel.Log{
			Action:   "chat_created",
			EntityID: chatID,
		}

		repoErr = errors.New("repository error")
		logErr  = errors.New("log error")
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: chatID,
			err:  nil,
			chatRepositoryMock: func(mc *minimock.Controller) *mocks.ChatRepositoryMock {
				mock := mocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chatModel).Return(chatID, nil)
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
			chatRepositoryMock: func(mc *minimock.Controller) *mocks.ChatRepositoryMock {
				mock := mocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chatModel).Return(0, repoErr)
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
			chatRepositoryMock: func(mc *minimock.Controller) *mocks.ChatRepositoryMock {
				mock := mocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, chatModel).Return(chatID, nil)
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
			chatRepoMock := tt.chatRepositoryMock(mc)
			logRepoMock := tt.logRepositoryMock(mc)
			
			txManager := &txManagerMock{}

			service := chatService.NewService(
				chatRepoMock,
				logRepoMock,
				txManager,
			)

			api := chat.NewImplementation(service)

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
