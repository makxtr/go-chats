package chat

import (
	"chat-server/internal/converter"
	desc "chat-server/pkg/chat_server_v1"
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	err := i.chatService.SendMessage(ctx, converter.ToMessageFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
