package log

import (
	"auth/internal/model"
	"auth/internal/repository"
	"context"
	"log"

	"github.com/makxtr/go-common/pkg/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "user_logs"

	actionColumn   = "action"
	entityIDColumn = "entity_id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.LogRepository {
	return &repo{db: db}
}

func (r *repo) Log(ctx context.Context, userLog *model.UserLog) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(actionColumn, entityIDColumn).
		Values(userLog.Action, userLog.EntityID)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	_, err = r.db.DB().ExecContext(ctx, db.Query{Name: "log_repository.Log", QueryRaw: query}, args...)
	if err != nil {
		log.Printf("failed to insert log: %v", err)
		return err
	}

	return nil
}
