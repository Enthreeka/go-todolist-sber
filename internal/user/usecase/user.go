package usecase

import (
	"context"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/entity"
	"go-todolist-sber/internal/user"
)

type userUsecase struct {
	userRepo user.UserRepository
	argon    *argon
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		argon:    NewArgonPassword(""),
	}
}

func (u *userUsecase) Login(ctx context.Context, user *entity.User) (*entity.User, error) {
	data, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		return nil, err
	}

	err = u.argon.VerifyPassword(data.Password, user.Password)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (u *userUsecase) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	if !entity.IsLoginValid(user.Login) {
		return nil, apperror.ErrDataNotValid
	}

	hashPassword, err := u.argon.GenerateHashFromPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashPassword
	data, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return data, nil
}
