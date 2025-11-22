package model

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Name  string
	Email string
	Role  Role
}

type CreateUserCommand struct {
	Info            UserInfo
	Password        string
	PasswordConfirm string
}

func (c *CreateUserCommand) Validate() error {
	if c.Password != c.PasswordConfirm {
		return errors.New("passwords do not match")
	}
	if len(c.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

type CreateUserData struct {
	Info           UserInfo
	HashedPassword string
}

type UpdateUserData struct {
	Name  *string
	Email *string
}

type Role int32

const (
	RoleUnspecified Role = iota
	RoleUser
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "USER"
	case RoleAdmin:
		return "ADMIN"
	default:
		return "UNSPECIFIED"
	}
}
