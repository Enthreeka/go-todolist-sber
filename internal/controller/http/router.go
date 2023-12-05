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

	auth := handler.AuthMiddleware(service.Session, store)

	mux.Route("/", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {

		})
		r.With(auth).Route("/task", func(r chi.Router) {
			r.Get("/list", task.GetTaskHandler)
			r.Post("/add", task.CreateTaskHandler)
			r.Delete("/delete", task.DeleteTaskHandler)
			r.Put("/update", task.UpdateTaskHandler)
			r.Get("/all", task.GetAllTasksHandler)
			r.Get("/pagination", task.GetTaskWithPaginationHandler)
			r.Get("/filter", task.GetFilteredHandler)
		})
	})

	return mux
}
