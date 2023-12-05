package handler

import (
	"context"
	"encoding/json"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/session"
	"go-todolist-sber/internal/user"
	"go-todolist-sber/pkg/logger"
	"net/http"
)

type userHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
	log            *logger.Logger
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase, log *logger.Logger) *userHandler {
	return &userHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
		log:            log,
	}
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (u *userHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := new(UserRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		u.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	user, err := u.userUsecase.Register(context.Background(), data.Login, data.Password)
	if err != nil {
		u.log.Error("taskUsecase.GetUserTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	_, err = u.sessionUsecase.CreateToken(context.Background(), user.ID)
	if err != nil {
		u.log.Error("taskUsecase.GetUserTasks: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(user)
}

func (u *userHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	data := new(UserRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		u.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	user, err := u.userUsecase.Login(context.Background(), data.Login, data.Password)
	if err != nil {
		u.log.Error("userUsecase.Login: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.Encode(user)
}
