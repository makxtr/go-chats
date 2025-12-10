package auth_test

import (
	"context"
	"errors"
	"testing"

	authAPI "auth/internal/api/auth"
	"auth/internal/model"
	"auth/internal/repository/mocks"
	authService "auth/internal/service/auth"
	desc "auth/pkg/auth_v1"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestImplementation_Login(t *testing.T) {
	type userRepositoryMockFunc func(mc *minimock.Controller) *mocks.UserRepositoryMock
	type logRepositoryMockFunc func(mc *minimock.Controller) *mocks.LogRepositoryMock

	type args struct {
		ctx context.Context
		req *desc.LoginRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username = "test_user"
		password = "password123"
		role     = model.RoleUser

		// Hash password for test
		hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		req = &desc.LoginRequest{
			Username: username,
			Password: password,
		}

		userSecure = &model.UserSecure{
			Password: string(hashedPassword),
			Role:     role,
		}

		repoErr = errors.New("repository error")
	)

	tests := []struct {
		name               string
		args               args
		wantErr            bool
		userRepositoryMock userRepositoryMockFunc
		logRepositoryMock  logRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			wantErr: false,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetForLoginMock.Expect(ctx, username).Return(userSecure, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				return mocks.NewLogRepositoryMock(mc)
			},
		},
		{
			name: "user not found",
			args: args{
				ctx: ctx,
				req: req,
			},
			wantErr: true,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetForLoginMock.Expect(ctx, username).Return(nil, repoErr)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				return mocks.NewLogRepositoryMock(mc)
			},
		},
		{
			name: "wrong password",
			args: args{
				ctx: ctx,
				req: &desc.LoginRequest{
					Username: username,
					Password: "wrong_password",
				},
			},
			wantErr: true,
			userRepositoryMock: func(mc *minimock.Controller) *mocks.UserRepositoryMock {
				mock := mocks.NewUserRepositoryMock(mc)
				mock.GetForLoginMock.Expect(ctx, username).Return(userSecure, nil)
				return mock
			},
			logRepositoryMock: func(mc *minimock.Controller) *mocks.LogRepositoryMock {
				return mocks.NewLogRepositoryMock(mc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userRepoMock := tt.userRepositoryMock(mc)
			logRepoMock := tt.logRepositoryMock(mc)

			txManager := &txManagerMock{}
			securityConfig := &securityConfigMock{
				refreshKey: "test-refresh-secret-key-32-bytes-long",
				accessKey:  "test-access-secret-key-32-bytes-long!",
				refreshExp: 60,
				accessExp:  5,
			}

			service := authService.NewService(
				userRepoMock,
				logRepoMock,
				txManager,
				securityConfig,
			)

			api := authAPI.NewImplementation(service)

			resp, err := api.Login(tt.args.ctx, tt.args.req)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.NotEmpty(t, resp.RefreshToken)
			}
		})
	}
}
