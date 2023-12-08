package user

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type UserUsecase interface {
	Login(ctx context.Context, login, password string) (*entity.User, error)
	Register(ctx context.Context, userID string, login, password string) (*entity.User, error)
}
