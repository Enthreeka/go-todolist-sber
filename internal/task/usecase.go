package task

import (
	"context"
	"go-todolist-sber/internal/entity"
	"time"
)

type TaskUsecase interface {
	GetFilteredTasks(ctx context.Context, userID string, date time.Time, status bool) ([]entity.Task, error)
	PaginationTasks(ctx context.Context, userID string, status bool, page int) ([]entity.Task, error)
	CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	GetUserTasks(ctx context.Context, userID string) ([]entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
}
