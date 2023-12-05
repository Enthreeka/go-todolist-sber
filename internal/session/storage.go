package session

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type SessionRepository interface {
	Create(ctx context.Context, session *entity.Session) (*entity.Session, error)
	GetByToken(ctx context.Context, token string) (*entity.Session, error)
	Update(ctx context.Context, session *entity.Session) (*entity.Session, error)
}
