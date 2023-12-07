package task

import (
	"context"
	"go-todolist-sber/internal/entity"
)

type TaskUsecase interface {
	GetTaskWithPaginationAndFilter(ctx context.Context, userID string, option *entity.ParamOption) ([]entity.Task, error)
	CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	DeleteTask(ctx context.Context, id int) error
	GetUserTasks(ctx context.Context, userID string, option *entity.ParamOption) ([]entity.Task, error)
	GetAllTasks(ctx context.Context) ([]entity.Task, error)
	IsEqualUserID(ctx context.Context, contextUserID string, taskID int) (bool, error)
	UpdateTaskStatus(ctx context.Context, status bool, taskID int) (*entity.Task, error)
}
