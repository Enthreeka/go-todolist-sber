package handler

import (
	"go-todolist-sber/internal/user"
	"go-todolist-sber/pkg/logger"
	"net/http"
)

type userHandler struct {
	userUsecase user.UserUsecase
	log         *logger.Logger
}

func NewUserHandler(userUsecase user.UserUsecase, log *logger.Logger) *userHandler {
	return &userHandler{
		userUsecase: userUsecase,
		log:         log,
	}
}

func (u *userHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

}

func (u *userHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {

}
