package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/task/mock"
	"testing"
	"time"
)

func TestTaskUsecase_CreateTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	userID := uuid.New().String()

	task := &entity.Task{
		ID:          1,
		Done:        false,
		UserID:      userID,
		Header:      "Header",
		Description: "Description",
		CreatedAt:   time.Now(),
		StartDate:   time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()

	mockTaskRepo.EXPECT().Create(ctx, gomock.Eq(task)).Return(task, nil)

	createdTask, err := taskUsecase.CreateTask(ctx, task)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, createdTask)
}

func TestTaskUsecase_DeleteTask(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *mock.MockTaskRepository, id int)
	type args struct {
		id int
	}

	tests := []struct {
		name string
		args
		mockBehavior
		want    string
		wantErr error
	}{
		{
			name: "ok",
			args: args{id: 6},
			mockBehavior: func(m *mock.MockTaskRepository, id int) {
				m.EXPECT().DeleteByID(context.Background(), gomock.Eq(id)).Return(nil)
			},
			want:    "",
			wantErr: nil,
		},
		{
			name: "Not Found",
			args: args{id: 15},
			mockBehavior: func(m *mock.MockTaskRepository, id int) {
				m.EXPECT().DeleteByID(context.Background(), gomock.Eq(id)).Return(apperror.ErrNoRows)
			},
			want:    "",
			wantErr: apperror.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTaskRepo := mock.NewMockTaskRepository(ctrl)
			tt.mockBehavior(mockTaskRepo, tt.id)
			taskUsecase := NewTaskUsecase(mockTaskRepo)
			err := taskUsecase.DeleteTask(context.Background(), tt.args.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, "")
		})
	}
}

func TestTaskUsecase_UpdateTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	userID := uuid.New().String()

	task := &entity.Task{
		ID:          1,
		UserID:      userID,
		Header:      "Update Header",
		Description: "Description",
		CreatedAt:   time.Now(),
	}

	newTask := &entity.Task{
		ID:          1,
		UserID:      userID,
		Header:      "Update Header",
		Description: "Description",
		CreatedAt:   time.Now(),
		StartDate:   time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()

	mockTaskRepo.EXPECT().Update(ctx, gomock.Eq(task)).Return(newTask, nil)

	updatedTask, err := taskUsecase.UpdateTask(ctx, task)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, updatedTask)
}

func TestTaskUsecase_GetUserTasks(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	userID := uuid.New().String()
	option := &entity.ParamOption{}

	mockTaskRepo.EXPECT().GetByUserID(gomock.Any(), gomock.Eq(userID)).Return([]entity.Task{}, nil)

	tasks, err := taskUsecase.GetUserTasks(context.Background(), userID, option)
	require.NoError(t, err)
	require.NotNil(t, tasks)
}

func TestTaskUsecase_GetAllTasks(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	mockTaskRepo.EXPECT().GetAll(gomock.Any()).Return([]entity.Task{}, nil)

	tasks, err := taskUsecase.GetAllTasks(context.Background())
	require.NoError(t, err)
	require.NotNil(t, tasks)
}

func TestTaskUsecase_IsEqualUserID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	contextUserID := "user1"
	taskID := 1

	mockTaskRepo.EXPECT().GetByID(gomock.Any(), gomock.Eq(taskID)).Return(&entity.Task{UserID: contextUserID}, nil)

	result, err := taskUsecase.IsEqualUserID(context.Background(), contextUserID, taskID)
	require.NoError(t, err)
	require.True(t, result)
}

func TestTaskUsecase_UpdateTaskStatus(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	taskID := 1
	status := true

	mockTaskRepo.EXPECT().UpdateDone(gomock.Any(), gomock.Eq(status), gomock.Eq(taskID)).Return(&entity.Task{}, nil)

	updatedTask, err := taskUsecase.UpdateTaskStatus(context.Background(), status, taskID)
	require.NoError(t, err)
	require.NotNil(t, updatedTask)
}

func TestTaskUsecase_GetTask(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskRepo := mock.NewMockTaskRepository(ctrl)

	taskUsecase := NewTaskUsecase(mockTaskRepo)

	userID := uuid.New().String()
	option := &entity.ParamOption{
		Page: 1,
	}

	offset := (option.Page - 1) * 3
	mockTaskRepo.EXPECT().GetByUserIDWithOffset(gomock.Any(), gomock.Eq(userID), gomock.Eq(offset)).Return([]entity.Task{}, nil)

	tasks, err := taskUsecase.GetTask(context.Background(), userID, option)
	require.NoError(t, err)
	require.NotNil(t, tasks)
}
