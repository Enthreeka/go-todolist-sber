package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go-todolist-sber/pkg/logger"
	"net/http"
)

type ServerOption struct {
	Addr string
}

func NewServer(log *logger.Logger, services Services, opts ServerOption) *http.Server {
	mux := chi.NewMux()

	mux.Use(middleware.RealIP,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
	)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	return &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}
}