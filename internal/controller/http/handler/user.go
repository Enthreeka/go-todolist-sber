package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"go-todolist-sber/internal/apperror"
	"go-todolist-sber/internal/session"
	"go-todolist-sber/internal/user"
	"go-todolist-sber/pkg/logger"
	"net/http"
)

type userHandler struct {
	userUsecase    user.UserUsecase
	sessionUsecase session.SessionUsecase
	store          *sessions.CookieStore
	log            *logger.Logger
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase, store *sessions.CookieStore, log *logger.Logger) *userHandler {
	return &userHandler{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
		store:          store,
		log:            log,
	}
}

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// RegisterHandler godoc
// @Summary Register new user
// @Tags Auth
// @Description register new user, returns user and set session
// @Accept json
// @Produce json
// @Param input body UserRequest true "user login and password"
// @Success 201 {object} entity.User
// @Failure 400 {object} JSONError
// @Failure 422 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /user/register [post]
func (u *userHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	data := new(UserRequest)
	d := json.NewDecoder(r.Body)
	err := d.Decode(&data)
	if err != nil {
		u.log.Error("json.NewDecoder: %v", err)
		DecodingError(w)
		return
	}

	user, err := u.userUsecase.Register(context.Background(), uuid.New().String(), data.Login, data.Password)
	if err != nil {
		u.log.Error("userUsecase.Register: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	sess, err := u.sessionUsecase.CreateToken(context.Background(), user.ID)
	if err != nil {
		u.log.Error("sessionUsecase.CreateToken: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	u.authenticated(w, r, sess.Token, true)

	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(user)
}

// LoginHandler godoc
// @Summary Login user
// @Tags Auth
// @Description login user,returns user and set session
// @Accept json
// @Produce json
// @Param input body UserRequest true "user login and password"
// @Success 200 {object} entity.User
// @Failure 400 {object} JSONError
// @Failure 401 {object} JSONError
// @Failure 404 {object} JSONError
// @Failure 500 {object} JSONError
// @Router /user/login [post]
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

	sess, err := u.sessionUsecase.UpdateToken(context.Background(), user.ID)
	if err != nil {
		u.log.Error("sessionUsecase.UpdateToken: %v", err)
		HandleError(w, err, apperror.ParseHTTPErrStatusCode(err))
		return
	}

	u.authenticated(w, r, sess.Token, true)

	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent(" ", " ")
	e.Encode(user)
}

// LogoutHandler godoc
// @Summary Logout user
// @Tags Auth
// @Description logout user removing session
// @Accept json
// @Produce json
// @Success 200
// @Router /user/logout [post]
func (u *userHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	u.authenticated(w, r, "", false)
	w.WriteHeader(http.StatusOK)
}

func (u *userHandler) authenticated(w http.ResponseWriter, r *http.Request, sessionID string, authenticated bool) {
	session, err := u.store.Get(r, "session.id")
	if err != nil {
		u.log.Error("store.Get(r,sessionID): %v", err)
	}
	session.Values["authenticated"] = authenticated
	session.Values["sessionID"] = sessionID

	session.Save(r, w)
}
