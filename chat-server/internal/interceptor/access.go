package interceptor

import (
	"context"
	"log"
	"strings"

	accesspb "chat-server/pkg/access_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	authServiceAddr string
	accessClient    accesspb.AccessV1Client
}

func NewAuthInterceptor(authServiceAddr string) (*AuthInterceptor, error) {
	var opts []grpc.DialOption

	if strings.HasSuffix(authServiceAddr, ":443") {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(
		authServiceAddr,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	accessClient := accesspb.NewAccessV1Client(conn)

	return &AuthInterceptor{
		authServiceAddr: authServiceAddr,
		accessClient:    accessClient,
	}, nil
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("Auth check for method: %s", info.FullMethod)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not found")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token required")
		}

		outgoingCtx := metadata.NewOutgoingContext(ctx, md)

		checkReq := &accesspb.CheckRequest{
			EndpointAddress: info.FullMethod,
		}

		_, err := i.accessClient.Check(outgoingCtx, checkReq)
		if err != nil {
			log.Printf("Auth check failed: %v", err)
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		log.Printf("Auth check passed for %s", info.FullMethod)

		return handler(ctx, req)
	}
}
