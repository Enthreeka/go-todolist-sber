package usecase

import (
	"context"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/user"
)

type userUsecase struct {
	userRepo user.UserRepository
	argon    Argon
}

func NewUserUsecase(userRepo user.UserRepository, salt string) user.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		argon:    NewArgonPassword(salt),
	}
}

func (u *userUsecase) Login(ctx context.Context, login, password string) (*entity.User, error) {
	if !entity.IsLoginValid(login) {
		return nil, apperror.ErrDataNotValid
	}
	if !entity.IsPasswordValid(password) {
		return nil, apperror.ErrDataNotValid
	}

	data, err := u.userRepo.GetByLogin(ctx, login)
	if err != nil {
		return nil, err
	}

	err = u.argon.VerifyPassword(data.Password, password)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *userUsecase) Register(ctx context.Context, userID string, login, password string) (*entity.User, error) {
	if !entity.IsLoginValid(login) {
		return nil, apperror.ErrDataNotValid
	}

	hashPassword, err := u.argon.GenerateHashFromPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       userID,
		Password: hashPassword,
		Login:    login,
	}

	data, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return data, nil
}
