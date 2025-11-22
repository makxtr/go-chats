package model

type Chat struct {
	ID        int64    `db:"id"`
	Usernames []string `db:"usernames"`
}
