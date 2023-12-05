package server

import (
	"context"
	"fmt"
	"go-todolist-sber/internal/config"
	"go-todolist-sber/internal/controller/http"
	"go-todolist-sber/internal/task/repo"
	"go-todolist-sber/internal/task/usecase"
	"go-todolist-sber/pkg/logger"
	"go-todolist-sber/pkg/postgres"
)

func Run(log *logger.Logger, cfg *config.Config) error {
	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}

	defer psql.Close()

	taskRepo := repo.NewTaskRepository(psql)

	taskUsecase := usecase.NewTaskUsecase(taskRepo)

	server := http.NewServer(log, http.Services{Task: taskUsecase}, http.ServerOption{
		Addr: fmt.Sprintf("%s:%s", cfg.HTTTPServer.Hostname, cfg.HTTTPServer.Port),
	})

	log.Info("Starting http server on %s: %s%s", cfg.HTTTPServer.TypeServer, cfg.HTTTPServer.Hostname, cfg.HTTTPServer.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("listen: %s", err)
	}

	return nil
}
