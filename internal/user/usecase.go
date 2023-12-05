package user

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type UserUsecase interface {
	Login(ctx context.Context, user *entity.User) (*entity.User, error)
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
}
