package repository

import "errors"

var (
	ErrNotFound      = errors.New("entity not found")
	ErrQueryBuild    = errors.New("failed to build query")
	ErrQueryExec     = errors.New("failed to execute query")
	ErrCreateFailed  = errors.New("failed to create entity")
	ErrUpdateFailed  = errors.New("failed to update entity")
	ErrDeleteFailed  = errors.New("failed to delete entity")
)