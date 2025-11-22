package chat

import (
	"chat-server/internal/model"
	"context"
	"log"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	// In the future, we might want to save the message to the DB or send it to a message broker.
	// For now, we just log it as per the original implementation.
	log.Printf("Message received - From: %s, Text: %s, Timestamp: %v",
		message.From,
		message.Text,
		message.Timestamp)

	return nil
}
