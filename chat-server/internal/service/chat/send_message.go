package chat

import (
	"chat-server/internal/model"
	"context"
	"log"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		// In the future, we might want to save the message to the DB or send it to a message broker.
		// For now, we just log it as per the original implementation.
		log.Printf("Message received - From: %s, Text: %s, Timestamp: %v",
			message.From,
			message.Text,
			message.Timestamp)
		// Since the current API definition does not include a ChatID in SendMessageRequest,
		// we cannot log a meaningful EntityID. We will use 0 for now.
		errTx = s.logRepository.Log(ctx, &model.ChatLog{
			Action:   "message_sent",
			EntityID: 0,
		})
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
