package http

import (
	"github.com/go-chi/chi/v5"
	"go-todolist-sber/internal/controller/http/handler"
	"go-todolist-sber/internal/usecase"
)

type Services struct {
	Task usecase.Task
}

func Router(service Services) *chi.Mux {
	mux := chi.NewMux()

	task := handler.NewTaskHandler(service.Task)

	mux.Route("/", func(r chi.Router) {
		r.Route("/task", func(r chi.Router) {
			r.Get("/list", task.GetTaskHandler)
		})
	})

	return mux
}
