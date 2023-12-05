package handler

import (
	"github.com/gorilla/sessions"
	"go-todolist-sber/pkg/logger"
	"net/http"
	"time"
)

type middleware func(next http.Handler) http.Handler

func MiddlewareLogger(log *logger.Logger) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			log.Info("%s - %s (%v)", r.Method, r.URL.Path, time.Since(startTime))
			next.ServeHTTP(w, r)
		})
	}
}

func AuthMiddleware(store *sessions.CookieStore) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := store.Get(r, "session.id")
			if err != nil {
				ErrorJSON(w, "Internal server error", http.StatusUnauthorized)
			}

			next.ServeHTTP(w, r)
		})
	}
}
