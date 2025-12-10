package access_test

import (
	"context"
	"testing"
	"time"

	accessAPI "auth/internal/api/access"
	"auth/internal/model"
	accessService "auth/internal/service/access"
	"auth/internal/utils"
	desc "auth/pkg/access_v1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestImplementation_Check(t *testing.T) {
	type args struct {
		ctx context.Context
		req *desc.CheckRequest
	}

	var (
		secretKey = "test-access-secret-key-32-bytes-long!"

		// Generate valid token
		validToken, _ = utils.GenerateToken(
			model.UserInfo{
				Name: "test_user",
				Role: model.RoleAdmin,
			},
			[]byte(secretKey),
			time.Hour,
		)

		// Generate token with user role
		userRoleToken, _ = utils.GenerateToken(
			model.UserInfo{
				Name: "test_user",
				Role: model.RoleUser,
			},
			[]byte(secretKey),
			time.Hour,
		)

		// Generate expired token
		expiredToken, _ = utils.GenerateToken(
			model.UserInfo{
				Name: "test_user",
				Role: model.RoleAdmin,
			},
			[]byte(secretKey),
			-time.Hour, // Already expired
		)

		protectedEndpoint = model.ExamplePath       // Requires admin role
		publicEndpoint    = "/some/public/endpoint" // Not in access map
	)

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		errCode    codes.Code
		errMessage string
	}{
		{
			name: "success case - valid admin token for protected endpoint",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("authorization", "Bearer "+validToken),
				),
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint,
				},
			},
			wantErr: false,
		},
		{
			name: "success case - public endpoint (not in access map)",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("authorization", "Bearer "+validToken),
				),
				req: &desc.CheckRequest{
					EndpointAddress: publicEndpoint,
				},
			},
			wantErr: false,
		},
		{
			name: "error - no metadata in context",
			args: args{
				ctx: context.Background(), // No metadata
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint,
				},
			},
			wantErr:    true,
			errCode:    codes.Unauthenticated,
			errMessage: "Authorization token is missing or invalid format",
		},
		{
			name: "error - no authorization header",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("some-other-header", "value"),
				),
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint,
				},
			},
			wantErr:    true,
			errCode:    codes.Unauthenticated,
			errMessage: "Authorization token is missing or invalid format",
		},
		{
			name: "error - invalid token format (no Bearer prefix)",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("authorization", validToken), // Missing "Bearer "
				),
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint,
				},
			},
			wantErr:    true,
			errCode:    codes.Unauthenticated,
			errMessage: "Authorization token is missing or invalid format",
		},
		{
			name: "error - expired token",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("authorization", "Bearer "+expiredToken),
				),
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint,
				},
			},
			wantErr:    true,
			errCode:    codes.Unauthenticated,
			errMessage: "Token verification failed (expired or invalid signature)",
		},
		{
			name: "error - invalid token signature",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("authorization", "Bearer invalid.jwt.token"),
				),
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint,
				},
			},
			wantErr:    true,
			errCode:    codes.Unauthenticated,
			errMessage: "Token verification failed (expired or invalid signature)",
		},
		{
			name: "error - access denied (user role trying to access admin endpoint)",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(),
					metadata.Pairs("authorization", "Bearer "+userRoleToken),
				),
				req: &desc.CheckRequest{
					EndpointAddress: protectedEndpoint, // Requires admin
				},
			},
			wantErr:    true,
			errCode:    codes.PermissionDenied,
			errMessage: "Access denied for this role",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			securityConfig := &securityConfigMock{
				accessKey: secretKey,
			}

			service := accessService.NewService(securityConfig)
			api := accessAPI.NewImplementation(service)

			resp, err := api.Check(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)

				// Check gRPC status code
				st, ok := status.FromError(err)
				require.True(t, ok, "error should be a gRPC status error")
				require.Equal(t, tt.errCode, st.Code())
				require.Contains(t, st.Message(), tt.errMessage)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
			}
		})
	}
}
