package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/session/mock"
	"testing"
	"time"
)

func TestSessionUsecase_CreateToken(t *testing.T) {
	type mockBehavior func(r *mock.MockSessionRepository, session *entity.Session)

	type args struct {
		session *entity.Session
		userID  string
		token   string
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *entity.Session
		wantErr      error
	}{
		{
			name: "ok",
			args: args{session: &entity.Session{
				UserID:    "uuid",
				Token:     "token",
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
				token:  "token",
				userID: "uuid"},
			mockBehavior: func(r *mock.MockSessionRepository, session *entity.Session) {
				r.EXPECT().Create(context.Background(), session).Return(&entity.Session{}, nil)
			},
			want:    &entity.Session{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSessionRepo := mock.NewMockSessionRepository(ctrl)
			tt.mockBehavior(mockSessionRepo, tt.args.session)

			sessionUsecase := NewSessionUsecase(mockSessionRepo)

			sess, err := sessionUsecase.CreateToken(context.Background(), tt.args.token, tt.args.userID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, sess)
		})
	}
}

func TestSessionUsecase_GetToken(t *testing.T) {
	type mockBehavior func(r *mock.MockSessionRepository, token string)

	type args struct {
		token string
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *entity.Session
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				token: uuid.New().String(),
			},
			mockBehavior: func(r *mock.MockSessionRepository, token string) {
				r.EXPECT().GetByToken(context.Background(), token).Return(&entity.Session{}, nil)
			},
			want:    &entity.Session{},
			wantErr: nil,
		},
		{
			name: "not found",
			args: args{
				token: uuid.New().String(),
			},
			mockBehavior: func(r *mock.MockSessionRepository, token string) {
				r.EXPECT().GetByToken(context.Background(), token).Return(nil, apperror.ErrNoRows)
			},
			want:    nil,
			wantErr: apperror.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSessionRepo := mock.NewMockSessionRepository(ctrl)
			tt.mockBehavior(mockSessionRepo, tt.args.token)

			sessionUsecase := NewSessionUsecase(mockSessionRepo)

			sess, err := sessionUsecase.GetToken(context.Background(), tt.args.token)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, sess)
		})
	}
}

func TestSessionUsecase_UpdateToken(t *testing.T) {
	type mockBehavior func(r *mock.MockSessionRepository, session *entity.Session)

	type args struct {
		session *entity.Session
		userID  string
		token   string
	}

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         *entity.Session
		wantErr      error
	}{
		{
			name: "ok",
			args: args{session: &entity.Session{
				UserID:    "uuid",
				Token:     "token",
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
				token:  "token",
				userID: "uuid"},
			mockBehavior: func(r *mock.MockSessionRepository, session *entity.Session) {
				r.EXPECT().Update(context.Background(), session).Return(&entity.Session{}, nil)
			},
			want:    &entity.Session{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockSessionRepo := mock.NewMockSessionRepository(ctrl)
			tt.mockBehavior(mockSessionRepo, tt.args.session)

			sessionUsecase := NewSessionUsecase(mockSessionRepo)

			sess, err := sessionUsecase.UpdateToken(context.Background(), tt.args.token, tt.args.userID)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, sess)
		})
	}
}
