package usecase

import (
	"context"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/task"
	"time"
)

type taskUsecase struct {
	taskRepo task.TaskRepository
}

func NewTaskUsecase(taskRepo task.TaskRepository) task.TaskUsecase {
	return &taskUsecase{
		taskRepo: taskRepo,
	}
}

func (t *taskUsecase) CreateTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	task, err := t.taskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *taskUsecase) UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	task, err := t.taskRepo.Update(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *taskUsecase) DeleteTask(ctx context.Context, id int) error {
	if err := t.taskRepo.DeleteByID(ctx, id); err != nil {
		return err
	}
	return nil
}

func (t *taskUsecase) GetUserTasks(ctx context.Context, userID string) ([]entity.Task, error) {
	tasks, err := t.taskRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if tasks == nil || len(tasks) == 0 {
		return nil, apperror.ErrNoRows
	}

	return tasks, nil
}

func (t *taskUsecase) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	tasks, err := t.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskUsecase) PaginationTasks(ctx context.Context, userID string, status bool, page int) ([]entity.Task, error) {
	offset := (page - 1) * 3

	tasks, err := t.taskRepo.GetPageByStatusAndUserID(ctx, userID, status, offset)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *taskUsecase) GetFilteredTasks(ctx context.Context, userID string, date time.Time, status bool) ([]entity.Task, error) {
	tasks, err := t.taskRepo.GetByDateAndStatus(ctx, userID, date, status)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}