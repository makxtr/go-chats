package chat_test

import (
	"context"

	"github.com/makxtr/go-common/pkg/db"
)

// txManagerMock is a simple mock for TxManager that executes the function without transaction
type txManagerMock struct{}

func (tm *txManagerMock) ReadCommitted(ctx context.Context, fn db.Handler) error {
	return fn(ctx)
}
