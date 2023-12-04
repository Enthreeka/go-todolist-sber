package main

import (
	"go-todolist-sber/internal/config"
	"go-todolist-sber/internal/server"
	"go-todolist-sber/pkg/logger"
)

func main() {
	log := logger.New()

	cfg, err := config.New()
	if err != nil {
		log.Fatal("failed load config: %v", err)
	}

	if err := server.Run(log, cfg); err != nil {
		log.Fatal("failed to run server: %v", err)
	}
}