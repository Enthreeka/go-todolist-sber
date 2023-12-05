package usecase

import (
	"go-todolist-sber/internal/user"
)

type userUsecase struct {
	userRepo user.UserRepository
}

func NewUserUsecase(userRepo user.UserRepository) user.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

//func (u *userUsecase) Login(ctx context.Context, user *entity.User) (*entity.User, error) {
//	data, err := u.userRepo.GetByLogin(ctx, user.Login)
//	if err != nil {
//		return nil, err
//	}
//
//}
