package chat

import (
	"chat-server/internal/converter"
	desc "chat-server/pkg/chat_server_v1"
	"context"
	"log"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, converter.ToChatFromDesc(req))
	if err != nil {
		return nil, err
	}

	log.Printf("created chat with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
