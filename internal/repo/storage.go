package repo

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type Task interface {
	GetPageByDoneAndUserID(ctx context.Context, userID string, done bool, offset int) ([]entity.Task, error)
	Create(ctx context.Context, task *entity.Task) (*entity.Task, error)
	GetByUserID(ctx context.Context, id string) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) (*entity.Task, error)
	GetAll(ctx context.Context) ([]entity.Task, error)
	DeleteByID(ctx context.Context, id int) error
}

type User interface {
}
