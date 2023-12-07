package task

import (
	"context"
	"go-todolist-sber/internal/entity"
	"time"
)

type TaskRepository interface {
	GetPageByStatusAndUserID(ctx context.Context, userID string, status bool, offset int) ([]entity.Task, error)
	GetByDateAndStatus(ctx context.Context, userID string, date time.Time, status bool, offset int) ([]entity.Task, error)
	Create(ctx context.Context, task *entity.Task) (*entity.Task, error)
	GetByUserID(ctx context.Context, id string) ([]entity.Task, error)
	Update(ctx context.Context, task *entity.Task) (*entity.Task, error)
	GetAll(ctx context.Context) ([]entity.Task, error)
	DeleteByID(ctx context.Context, id int) error
	GetByID(ctx context.Context, id int) (*entity.Task, error)
	UpdateDone(ctx context.Context, status bool, taskID int) (*entity.Task, error)
}
