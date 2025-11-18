package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	desc "auth/pkg/user_v1"
)

func main() {
	// Cloud Run требует TLS
	creds := credentials.NewClientTLSFromCert(nil, "")

	conn, err := grpc.NewClient(
		"auth-service-rxpqkfxb3a-uc.a.run.app:443",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: 12})
	if err != nil {
		log.Fatalf("failed to get user: %v", err)
	}

	log.Printf("✅ Success! User from Cloud: %+v", r.GetUser())
}