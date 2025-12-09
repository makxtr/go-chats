package auth_test

import (
	"context"

	"github.com/makxtr/go-common/pkg/db"
)

// txManagerMock is a simple mock for TxManager that executes the function without transaction
type txManagerMock struct{}

func (tm *txManagerMock) ReadCommitted(ctx context.Context, fn db.Handler) error {
	return fn(ctx)
}

// securityConfigMock is a simple mock for SecurityConfig
type securityConfigMock struct {
	refreshKey string
	accessKey  string
	refreshExp byte
	accessExp  byte
}

func (sc *securityConfigMock) RefreshKey() string {
	return sc.refreshKey
}

func (sc *securityConfigMock) AccessKey() string {
	return sc.accessKey
}

func (sc *securityConfigMock) RefreshExp() byte {
	return sc.refreshExp
}

func (sc *securityConfigMock) AccessExp() byte {
	return sc.accessExp
}