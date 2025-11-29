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
	"time"

	"github.com/gojuno/minimock/v3"
	logModel "github.com/makxtr/go-common/pkg/logger/model"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestImplementation_SendMessage(t *testing.T) {
	type chatRepositoryMockFunc func(mc *minimock.Controller) *mocks.ChatRepositoryMock
	type logRepositoryMockFunc func(mc *minimock.Controller) *mocks.LogRepositoryMock

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		timestamp = time.Now()

		req = &desc.SendMessageRequest{
			Message: &desc.Message{
				From:      "user1",
				Text:      "Hello, World!",
				Timestamp: timestamppb.New(timestamp),
			},
		}

		messageModel = &model.Message{
			From:      "user1",
			Text:      "Hello, World!",
			Timestamp: timestamp,
		}

		logEntry = &logModel.Log{
			Action:   "message_sent",
			EntityID: 0,
		}

		logErr = errors.New("log error")
	)

	tests := []struct {
		name               string
		args               args
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
			err: nil,
			chatRepositoryMock: func(mc *minimock.Controller) *mocks.ChatRepositoryMock {
				return mocks.NewChatRepositoryMock(mc)
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				mock := mocks.NewLogRepositoryMock(mc)
				mock.LogMock.Expect(ctx, logEntry).Return(nil)
				return mock
			},
		},
		{
			name: "log error",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: logErr,
			chatRepositoryMock: func(mc *minimock.Controller) *mocks.ChatRepositoryMock {
				return mocks.NewChatRepositoryMock(mc)
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

			resp, err := api.SendMessage(tt.args.ctx, tt.args.req)

			if tt.err != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err.Error())
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}

			_ = messageModel
		})
	}
}
