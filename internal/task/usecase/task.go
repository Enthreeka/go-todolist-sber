package usecase

import (
	"context"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/task"
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

func (t *taskUsecase) GetUserTasks(ctx context.Context, userID string, option *entity.ParamOption) ([]entity.Task, error) {
	if option.Status != nil && !option.DateTime.IsZero() {
		tasks, err := t.taskRepo.GetByDateAndStatus(ctx, userID, option.DateTime, *option.Status)
		if err != nil {
			return nil, err
		}

		if tasks == nil || len(tasks) == 0 {
			return nil, apperror.ErrNoRows
		}

		return tasks, nil
	} else {

		tasks, err := t.taskRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if tasks == nil || len(tasks) == 0 {
			return nil, apperror.ErrNoRows
		}

		return tasks, nil
	}
}

func (t *taskUsecase) GetAllTasks(ctx context.Context) ([]entity.Task, error) {
	tasks, err := t.taskRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if tasks == nil || len(tasks) == 0 {
		return nil, apperror.ErrNoRows
	}

	return tasks, nil
}

func (t *taskUsecase) IsEqualUserID(ctx context.Context, contextUserID string, taskID int) (bool, error) {
	data, err := t.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return false, err
	}

	if data.UserID != contextUserID {
		return false, nil
	}

	return true, nil
}

func (t *taskUsecase) UpdateTaskStatus(ctx context.Context, status bool, taskID int) (*entity.Task, error) {
	task, err := t.taskRepo.UpdateDone(ctx, status, taskID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *taskUsecase) GetTaskWithPaginationAndFilter(ctx context.Context, userID string, option *entity.ParamOption) ([]entity.Task, error) {
	if option.Page == 0 {
		option.Page = 1
	}

	offset := (option.Page - 1) * 3
	if !option.DateTime.IsZero() && option.Status != nil {
		tasks, err := t.taskRepo.GetByDateAndStatusWithOffset(ctx, userID, option.DateTime, *option.Status, offset)
		if err != nil {
			return nil, err
		}

		if tasks == nil || len(tasks) == 0 {
			return nil, apperror.ErrNoRows
		}

		return tasks, nil
	} else if option.Status != nil {
		tasks, err := t.taskRepo.GetPageByStatusAndUserID(ctx, userID, *option.Status, offset)
		if err != nil {
			return nil, err
		}

		if tasks == nil || len(tasks) == 0 {
			return nil, apperror.ErrNoRows
		}

		return tasks, nil
	}

	tasks, err := t.taskRepo.GetByUserIDWithOffset(ctx, userID, offset)
	if err != nil {
		return nil, err
	}

	if tasks == nil || len(tasks) == 0 {
		return nil, apperror.ErrNoRows
	}

	return tasks, nil
}
