package server

import (
	"context"
	"fmt"
	"github.com/gorilla/sessions"
	"go-todolist-sber/internal/config"
	"go-todolist-sber/internal/controller/http"
	sessionRepo "go-todolist-sber/internal/session/repo"
	sessionUsecase "go-todolist-sber/internal/session/usecase"
	taskRepo "go-todolist-sber/internal/task/repo"
	taskUsecase "go-todolist-sber/internal/task/usecase"
	userRepo "go-todolist-sber/internal/user/repo"
	userUsecase "go-todolist-sber/internal/user/usecase"
	"go-todolist-sber/pkg/logger"
	"go-todolist-sber/pkg/postgres"
)

func Run(log *logger.Logger, cfg *config.Config) error {
	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}

	defer psql.Close()

	taskRepo := taskRepo.NewTaskRepository(psql)
	userRepo := userRepo.NewUserRepository(psql)
	sessionRepo := sessionRepo.NewSessionRepository(psql)

	taskUsecase := taskUsecase.NewTaskUsecase(taskRepo)
	userUsecase := userUsecase.NewUserUsecase(userRepo, cfg.Salt)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRepo)

	var store = sessions.NewCookieStore([]byte("secret-key"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}

	server := http.NewServer(log, http.Services{Task: taskUsecase, User: userUsecase, Session: sessionUsecase}, http.ServerOption{
		Addr: fmt.Sprintf(":%s", cfg.HTTTPServer.Port),
	}, store)

	log.Info("Starting http server on %s: %s:%s", cfg.HTTTPServer.TypeServer, cfg.HTTTPServer.Hostname, cfg.HTTTPServer.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("listen: %s", err)
	}

	return nil
}
