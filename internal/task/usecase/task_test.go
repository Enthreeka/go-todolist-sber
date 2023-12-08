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

	type mockBehavior func(r *mock.MockTaskRepository, task *entity.Task)
	type args struct {
		task *entity.Task
	}

	tests := []struct {
		name string
		args
		mockBehavior
		want    *entity.Task
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				task: &entity.Task{
					ID:          1,
					UserID:      "uuid",
					Header:      "Update Header",
					Description: "Description",
					CreatedAt:   time.Now(),
				},
			},
			mockBehavior: func(m *mock.MockTaskRepository, task *entity.Task) {
				m.EXPECT().Update(context.Background(), gomock.Eq(task)).Return(task, nil)
			},
			want:    &entity.Task{ID: 1, UserID: "uuid", Header: "Update Header", Description: "Description", CreatedAt: time.Now()},
			wantErr: nil,
		},
		{
			name: "Not found",
			args: args{
				task: &entity.Task{
					ID:          2,
					UserID:      uuid.New().String(),
					Header:      "Update Header 2",
					Description: "Description 2",
					CreatedAt:   time.Now(),
				},
			},
			mockBehavior: func(m *mock.MockTaskRepository, task *entity.Task) {
				m.EXPECT().Update(context.Background(), gomock.Eq(task)).Return(nil, apperror.ErrNoRows)
			},
			want:    nil,
			wantErr: apperror.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTaskRepo := mock.NewMockTaskRepository(ctrl)
			tt.mockBehavior(mockTaskRepo, tt.task)
			taskUsecase := NewTaskUsecase(mockTaskRepo)
			updatedTask, err := taskUsecase.UpdateTask(context.Background(), tt.task)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, updatedTask)
		})
	}
}

func TestTaskUsecase_IsEqualUserID(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *mock.MockTaskRepository, contextUserID string, taskID int)
	type args struct {
		contextUserID string
		taskID        int
	}
	tests := []struct {
		name string
		args
		mockBehavior
		want    bool
		wantErr error
	}{
		{
			name: "ok",
			args: args{
				contextUserID: "uuid",
				taskID:        5,
			},
			mockBehavior: func(m *mock.MockTaskRepository, contextUserID string, taskID int) {
				m.EXPECT().GetByID(context.Background(), gomock.Eq(taskID)).Return(&entity.Task{UserID: contextUserID}, nil)
			},
			want:    true,
			wantErr: nil,
		},
		{
			name: "Not found",
			args: args{
				contextUserID: "uuid",
				taskID:        5,
			},
			mockBehavior: func(m *mock.MockTaskRepository, contextUserID string, taskID int) {
				m.EXPECT().GetByID(context.Background(), gomock.Eq(taskID)).Return(nil, apperror.ErrNoRows)
			},
			want:    false,
			wantErr: apperror.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTaskRepo := mock.NewMockTaskRepository(ctrl)
			tt.mockBehavior(mockTaskRepo, tt.contextUserID, tt.taskID)
			taskUsecase := NewTaskUsecase(mockTaskRepo)
			equal, err := taskUsecase.IsEqualUserID(context.Background(), tt.contextUserID, tt.taskID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, equal)
		})
	}
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

	type mockBehavior func(r *mock.MockTaskRepository, userID string, date time.Time, status bool, offset int)
	type args struct {
		userID string
		option *entity.ParamOption
		offset int
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         []entity.Task
		wantErr      error
	}{
		{
			name: "ok with date and status",
			args: args{
				userID: "uuid",
				option: &entity.ParamOption{
					Page:     0,
					Status:   &[]bool{true}[0],
					DateTime: time.Now(),
				},
			},
			mockBehavior: func(m *mock.MockTaskRepository, userID string, date time.Time, status bool, offset int) {
				m.EXPECT().GetByDateAndStatus(context.Background(), gomock.Eq(userID), gomock.Eq(date), gomock.Eq(status)).Return([]entity.Task{}, nil)
			},
			want:    []entity.Task{},
			wantErr: nil,
		},
		{
			name: "ok with pagination and filter",
			args: args{
				userID: "uuid",
				option: &entity.ParamOption{
					Page:     1,
					Status:   &[]bool{true}[0],
					DateTime: time.Now(),
				},
				offset: 0,
			},
			mockBehavior: func(m *mock.MockTaskRepository, userID string, date time.Time, status bool, offset int) {
				m.EXPECT().GetByDateAndStatusWithOffset(context.Background(), userID, date, status, offset).Return([]entity.Task{}, nil)
			},
			want:    []entity.Task{},
			wantErr: nil,
		},
		{
			name: "ok with pagination",
			args: args{
				userID: "uuid",
				option: &entity.ParamOption{
					Page:   3,
					Status: &[]bool{true}[0],
				},
				offset: 6,
			},
			mockBehavior: func(m *mock.MockTaskRepository, userID string, date time.Time, status bool, offset int) {
				m.EXPECT().GetByStatusWithOffset(context.Background(), userID, status, offset).Return([]entity.Task{}, nil)
			},
			want:    []entity.Task{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockTaskRepo := mock.NewMockTaskRepository(ctrl)
			tt.mockBehavior(mockTaskRepo, tt.args.userID, tt.args.option.DateTime, *tt.args.option.Status, tt.args.offset)

			taskUsecase := NewTaskUsecase(mockTaskRepo)
			tasks, err := taskUsecase.GetTask(context.Background(), tt.args.userID, tt.args.option)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, tasks)
		})
	}
}
