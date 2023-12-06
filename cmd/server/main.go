package main

import (
	_ "go-todolist-sber/docs"
	"go-todolist-sber/internal/config"
	"go-todolist-sber/internal/server"
	"go-todolist-sber/pkg/logger"
)

// @title Blueprint Swagger API
// @version 1.0
// @description Swagger API for to do list
// @host localhost:8080
// @BasePath /
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
