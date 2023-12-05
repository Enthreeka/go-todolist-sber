package http

import (
	"github.com/go-chi/chi/v5"
	"go-todolist-sber/internal/controller/http/handler"
	"go-todolist-sber/internal/task"
	"go-todolist-sber/pkg/logger"
)

type Services struct {
	Task task.TaskUsecase
}

func Router(log *logger.Logger, service Services) *chi.Mux {
	mux := chi.NewMux()

	task := handler.NewTaskHandler(service.Task, log)

	mux.Route("/", func(r chi.Router) {
		r.Route("/task", func(r chi.Router) {
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
