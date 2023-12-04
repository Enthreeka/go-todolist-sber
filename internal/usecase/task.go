package usecase

import (
	"context"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/repo"
)

type taskUsecase struct {
	taskRepo repo.Task
}

func NewTaskUsecase(taskRepo repo.Task) Task {
	return &taskUsecase{
		taskRepo: taskRepo,
	}
}

func (t *taskUsecase) CreateTask(ctx context.Context, task *entity.Task) error {
	if _, err := t.taskRepo.Create(ctx, task); err != nil {
		return err
	}
	return nil
}

func (t *taskUsecase) UpdateTask(ctx context.Context, task *entity.Task) error {
	if _, err := t.taskRepo.Update(ctx, task); err != nil {
		return err
	}
	return nil
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

	return tasks, nil
}

func (t *taskUsecase) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	tasks, err := t.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}