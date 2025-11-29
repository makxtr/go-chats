package chat

import (
	"chat-server/internal/model"
	"chat-server/internal/repository"
	"context"
	"log"

	"github.com/makxtr/go-common/pkg/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
)

const (
	tableName = "chats"

	idColumn        = "id"
	usernamesColumn = "usernames"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernamesColumn).
		Values(pq.Array(chat.Usernames)).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, err
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, db.Query{Name: "chat_repository.Create", QueryRaw: query}, args...).Scan(&id)
	if err != nil {
		log.Printf("failed to create chat: %v", err)
		return 0, err
	}

	return id, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	res, err := r.db.DB().ExecContext(ctx, db.Query{Name: "chat_repository.Delete", QueryRaw: query}, args...)
	if err != nil {
		log.Printf("failed to delete chat: %v", err)
		return err
	}

	if res.RowsAffected() == 0 {
		return nil // Or return an error if you want to be strict
	}

	return nil
}
