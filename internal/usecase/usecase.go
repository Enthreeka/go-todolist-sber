package usecase

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type Task interface {
	PaginationTasks(ctx context.Context, userID string, done bool, page int) ([]entity.Task, error)
	CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	GetUserTasks(ctx context.Context, userID string) ([]entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}

type User interface {
}
