package user

import (
	"context"
	"go-todolist-sber/internal/entity"
)

//go:generate mockgen -source storage.go -destination mock/pg_repository_mock.go -package mock
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
}
