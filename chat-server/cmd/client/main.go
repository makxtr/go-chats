package main

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "chat-server/pkg/chat_server_v1"
)

const (
	address = "localhost:50052"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewChatServerV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SendMessage(ctx, &desc.SendMessageRequest{
		Message: &desc.Message{
			From:      gofakeit.Name(),
			Text:      gofakeit.Sentence(5),
			Timestamp: timestamppb.New(gofakeit.Date()),
		},
	})
	if err != nil {
		log.Fatalf("failed to send message: %v", err)
	}

	log.Printf("%s\n", color.RedString("Message sent"))
}
