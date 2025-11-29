package user

import (
	"auth/internal/model"
	"auth/internal/repository"
	repoConverter "auth/internal/repository/user/converter"
	modelRepo "auth/internal/repository/user/model"
	"context"
	"log"
	"time"

	"github.com/makxtr/go-common/pkg/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passColumn      = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, createUser *model.CreateUserData) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password", "role").
		Values(createUser.Info.Name, createUser.Info.Email, createUser.HashedPassword, int32(createUser.Info.Role)).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return 0, repository.ErrQueryBuild
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, db.Query{Name: "user_repository.Create", QueryRaw: query}, args...).Scan(&userID)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		return 0, repository.ErrCreateFailed
	}

	return userID, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, repository.ErrQueryBuild
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.ID, &user.Info.Name, &user.Info.Email, &user.Info.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return repoConverter.ToUserFromRepo(&user), nil
}

func (r *repo) Update(ctx context.Context, id int64, updateUser *model.UpdateUserData) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: id})

	if updateUser.Name != nil {
		builder = builder.Set(nameColumn, updateUser.Name)
	}

	if updateUser.Email != nil {
		builder = builder.Set(emailColumn, updateUser.Email)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return repository.ErrQueryBuild
	}

	res, err := r.db.DB().ExecContext(ctx, db.Query{Name: "user_repository.Update", QueryRaw: query}, args...)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return repository.ErrUpdateFailed
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return repository.ErrQueryBuild
	}

	res, err := r.db.DB().ExecContext(ctx, db.Query{Name: "user_repository.Delete", QueryRaw: query}, args...)
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return repository.ErrDeleteFailed
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
