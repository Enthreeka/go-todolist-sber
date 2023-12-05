package user

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
}
