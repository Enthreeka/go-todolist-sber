package usecase

import (
	"context"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/session"
	"time"
)

type sessionUsecase struct {
	sessionRepo session.SessionRepository
}

func NewSessionUsecase(sessionRepo session.SessionRepository) session.SessionUsecase {
	return &sessionUsecase{
		sessionRepo: sessionRepo,
	}
}

func (s *sessionUsecase) CreateToken(ctx context.Context, token string, userID string) (*entity.Session, error) {
	session := &entity.Session{
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Token:     token,
		UserID:    userID,
	}

	data, err := s.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *sessionUsecase) GetToken(ctx context.Context, token string) (*entity.Session, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionUsecase) UpdateToken(ctx context.Context, token string, userID string) (*entity.Session, error) {
	session := &entity.Session{
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		UserID:    userID,
	}

	data, err := s.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, err
	}

	return data, nil
}
