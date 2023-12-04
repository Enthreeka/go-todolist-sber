package server

import (
	"context"
	"go-todolist-sber/internal/config"
	"go-todolist-sber/pkg/logger"
	"go-todolist-sber/pkg/postgres"
)

func Run(log *logger.Logger, cfg *config.Config) error {
	psql, err := postgres.New(context.Background(), 5, cfg.Postgres.URL)
	if err != nil {
		log.Fatal("failed to connect PostgreSQL: %v", err)
	}
	defer psql.Close()

	return nil
}
