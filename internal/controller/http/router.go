package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"go-todolist-sber/internal/controller/http/handler"
	"go-todolist-sber/internal/session"
	"go-todolist-sber/internal/task"
	"go-todolist-sber/internal/user"
	"go-todolist-sber/pkg/logger"
)

type Services struct {
	Task    task.TaskUsecase
	User    user.UserUsecase
	Session session.SessionUsecase
}

func Router(log *logger.Logger, service Services, store *sessions.CookieStore) *chi.Mux {
	mux := chi.NewMux()

	task := handler.NewTaskHandler(service.Task, log)
	user := handler.NewUserHandler(service.User, service.Session, store, log)

	auth := handler.AuthMiddleware(service.Session, store)

	mux.Route("/", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", user.RegisterHandler)
			r.Post("/login", user.LoginHandler)
			r.Post("/logout", user.LogoutHandler)
		})
		r.With(auth).Route("/task", func(r chi.Router) {
			r.Get("/list", task.GetTaskHandler)
			r.Post("/add", task.CreateTaskHandler)
			r.Delete("/{id}", task.DeleteTaskHandler)
			r.Put("/{id}", task.UpdateTaskHandler)
			r.Get("/all", task.GetAllTasksHandler)
			r.Get("/pagination", task.GetTaskWithPaginationHandler)
			r.Get("/filter", task.GetFilteredHandler)
		})
	})

	return mux
}
