package converter

import (
	"chat-server/internal/model"
	desc "chat-server/pkg/chat_server_v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToChatFromDesc(req *desc.CreateRequest) *model.Chat {
	return &model.Chat{
		Usernames: req.GetUsernames(),
	}
}

func ToMessageFromDesc(req *desc.SendMessageRequest) *model.Message {
	return &model.Message{
		From:      req.GetMessage().GetFrom(),
		Text:      req.GetMessage().GetText(),
		Timestamp: req.GetMessage().GetTimestamp().AsTime(),
	}
}

func ToDescFromMessage(message *model.Message) *desc.Message {
	return &desc.Message{
		From:      message.From,
		Text:      message.Text,
		Timestamp: timestamppb.New(message.Timestamp),
	}
}
