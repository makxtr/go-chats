package interceptor

import (
	"context"
	"log"

	accesspb "chat-server/pkg/access_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// AuthInterceptor handles authentication and authorization
type AuthInterceptor struct {
	authServiceAddr string
	accessClient    accesspb.AccessV1Client
}

// NewAuthInterceptor creates a new auth interceptor
func NewAuthInterceptor(authServiceAddr string) (*AuthInterceptor, error) {
	// Создаем connection к auth-service
	conn, err := grpc.NewClient(
		authServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// Создаем gRPC клиент для AccessV1
	accessClient := accesspb.NewAccessV1Client(conn)

	return &AuthInterceptor{
		authServiceAddr: authServiceAddr,
		accessClient:    accessClient,
	}, nil
}

// Unary returns unary server interceptor для проверки доступа
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Printf("Auth check for method: %s", info.FullMethod)

		// Извлекаем metadata из входящего контекста
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not found")
		}

		// Проверяем наличие authorization header
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token required")
		}

		// Создаем исходящий контекст с токеном для auth-service
		outgoingCtx := metadata.NewOutgoingContext(ctx, md)

		// Вызываем Check метод auth-service используя сгенерированный клиент
		checkReq := &accesspb.CheckRequest{
			EndpointAddress: info.FullMethod,
		}

		_, err := i.accessClient.Check(outgoingCtx, checkReq)
		if err != nil {
			log.Printf("Auth check failed: %v", err)
			return nil, status.Error(codes.PermissionDenied, "access denied")
		}

		log.Printf("Auth check passed for %s", info.FullMethod)

		// Продолжаем выполнение запроса
		return handler(ctx, req)
	}
}