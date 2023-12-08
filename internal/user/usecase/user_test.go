package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/user/mock"
	"testing"
)

func TestUserUsecase_Register(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *mock.MockUserRepository, user *entity.User)
	type args struct {
		user     *entity.User
		login    string
		password string
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *entity.User
		wantErr      error
	}{
		{
			name: "ok",
			args: args{user: &entity.User{
				ID:       uuid.New().String(),
				Password: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
				Login:    "Test12414",
			},
				login:    "Test12414",
				password: "Test12414",
			},
			mockBehavior: func(r *mock.MockUserRepository, user *entity.User) {
				r.EXPECT().Create(context.Background(), gomock.Eq(user)).Return(&entity.User{}, nil)
			},
			want:    &entity.User{},
			wantErr: nil,
		},
		{
			name: "login not valid",
			args: args{user: &entity.User{
				ID: uuid.New().String(),
			},
				login:    "Test",
				password: "Test12414",
			},
			mockBehavior: func(r *mock.MockUserRepository, user *entity.User) {},
			want:         nil,
			wantErr:      apperror.ErrDataNotValid,
		},
		{
			name: "password not valid",
			args: args{user: &entity.User{
				ID:       uuid.New().String(),
				Password: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
				Login:    "Test12414",
			},
				login:    "Test12414",
				password: "Test12414",
			},
			mockBehavior: func(r *mock.MockUserRepository, user *entity.User) {
				r.EXPECT().Create(context.Background(), gomock.Eq(user)).Return(nil, apperror.ErrUniqueViolation)
			},
			want:    nil,
			wantErr: apperror.ErrUniqueViolation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepo := mock.NewMockUserRepository(ctrl)
			tt.mockBehavior(mockUserRepo, tt.args.user)

			userUsecase := NewUserUsecase(mockUserRepo, "")
			createdUSer, err := userUsecase.Register(context.Background(), tt.args.user.ID, tt.args.login, tt.args.login)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, createdUSer)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	t.Parallel()

	type mockBehavior func(r *mock.MockUserRepository, user *entity.User, login string)
	type args struct {
		user     *entity.User
		login    string
		password string
	}
	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *entity.User
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				user: &entity.User{
					Login:    "sdgsdg123",
					Password: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
				},
				login:    "sdgsdg123",
				password: "Test12414",
			},
			mockBehavior: func(r *mock.MockUserRepository, user *entity.User, login string) {
				r.EXPECT().GetByLogin(context.Background(), gomock.Eq(login)).Return(user, nil)
			},
			want: &entity.User{
				Login:    "sdgsdg123",
				Password: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
			},
			wantErr: nil,
		},
		{
			name: "not valid",
			args: args{
				user: &entity.User{
					Login:    "sdgsdg123",
					Password: "516a85c22c382c8055366d34d3178fc0cd5b73b4e5d8185cb8d7a8a55c9e2141",
				},
				login:    "sd",
				password: "Te",
			},
			mockBehavior: func(r *mock.MockUserRepository, user *entity.User, login string) {},
			want:         nil,
			wantErr:      apperror.ErrDataNotValid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockUserRepo := mock.NewMockUserRepository(ctrl)
			tt.mockBehavior(mockUserRepo, tt.args.user, tt.args.login)

			userUsecase := NewUserUsecase(mockUserRepo, "")
			createdUSer, err := userUsecase.Login(context.Background(), tt.args.login, tt.args.password)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, createdUSer)
		})
	}
}
