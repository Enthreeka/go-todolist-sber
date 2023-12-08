package session

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type SessionUsecase interface {
	CreateToken(ctx context.Context, token string, userID string) (*entity.Session, error)
	UpdateToken(ctx context.Context, token string, userID string) (*entity.Session, error)
	GetToken(ctx context.Context, token string) (*entity.Session, error)
}
