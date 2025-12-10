package access_test

import (
	"context"
	"testing"
	"time"

	"auth/internal/model"
	accessService "auth/internal/service/access"
	"auth/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestService_Check(t *testing.T) {
	type args struct {
		ctx         context.Context
		accessToken string
		endPoint    string
	}

	var (
		ctx       = context.Background()
		secretKey = "test-access-secret-key-32-bytes-long!"

		// Generate valid admin token
		validAdminToken, _ = utils.GenerateToken(
			model.UserInfo{
				Name: "admin_user",
				Role: model.RoleAdmin,
			},
			[]byte(secretKey),
			time.Hour,
		)

		// Generate valid user token
		validUserToken, _ = utils.GenerateToken(
			model.UserInfo{
				Name: "regular_user",
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

		protectedEndpoint = model.ExamplePath        // Requires admin role
		publicEndpoint    = "/some/public/endpoint" // Not in access map
	)

	tests := []struct {
		name      string
		args      args
		want      bool
		wantErr   bool
		errString string
	}{
		{
			name: "success - admin token accessing protected endpoint",
			args: args{
				ctx:         ctx,
				accessToken: validAdminToken,
				endPoint:    protectedEndpoint,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "success - public endpoint (not in access map)",
			args: args{
				ctx:         ctx,
				accessToken: validAdminToken,
				endPoint:    publicEndpoint,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "success - user token accessing public endpoint",
			args: args{
				ctx:         ctx,
				accessToken: validUserToken,
				endPoint:    publicEndpoint,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "error - user token accessing admin-only endpoint",
			args: args{
				ctx:         ctx,
				accessToken: validUserToken,
				endPoint:    protectedEndpoint, // Requires admin
			},
			want:      false,
			wantErr:   true,
			errString: "access denied",
		},
		{
			name: "error - expired token",
			args: args{
				ctx:         ctx,
				accessToken: expiredToken,
				endPoint:    protectedEndpoint,
			},
			want:      false,
			wantErr:   true,
			errString: "access token is invalid",
		},
		{
			name: "error - malformed token",
			args: args{
				ctx:         ctx,
				accessToken: "invalid.jwt.token",
				endPoint:    protectedEndpoint,
			},
			want:      false,
			wantErr:   true,
			errString: "access token is invalid",
		},
		{
			name: "error - empty token",
			args: args{
				ctx:         ctx,
				accessToken: "",
				endPoint:    protectedEndpoint,
			},
			want:      false,
			wantErr:   true,
			errString: "access token is invalid",
		},
		{
			name: "error - token with wrong signature",
			args: args{
				ctx: ctx,
				accessToken: func() string {
					// Generate token with different secret
					token, _ := utils.GenerateToken(
						model.UserInfo{
							Name: "test_user",
							Role: model.RoleAdmin,
						},
						[]byte("wrong-secret-key-different-from-test"),
						time.Hour,
					)
					return token
				}(),
				endPoint: protectedEndpoint,
			},
			want:      false,
			wantErr:   true,
			errString: "access token is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			securityConfig := &securityConfigMock{
				accessKey: secretKey,
			}

			service := accessService.NewService(securityConfig)

			got, err := service.Check(tt.args.ctx, tt.args.accessToken, tt.args.endPoint)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errString)
				require.Equal(t, tt.want, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}