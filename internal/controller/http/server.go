package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "go-todolist-sber/docs"
	"go-todolist-sber/internal/controller/http/handler"
	"go-todolist-sber/pkg/logger"
	"net/http"
)

type ServerOption struct {
	Addr string
}

func NewServer(log *logger.Logger, services Services, opts ServerOption, store *sessions.CookieStore) *http.Server {
	mux := chi.NewMux()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	mux.Use(middleware.RealIP,
		middleware.Recoverer,
		middleware.AllowContentType("application/json"),
		handler.MiddlewareLogger(log),
	)

	mux.Mount("/", Router(log, services, store))
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}
}
