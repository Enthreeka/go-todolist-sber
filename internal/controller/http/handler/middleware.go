package handler

import (
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
