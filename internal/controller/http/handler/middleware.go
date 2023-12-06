package handler

import (
	"context"
	"github.com/gorilla/sessions"
	"go-todolist-sber/internal/session"
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

func AuthMiddleware(sess session.SessionUsecase, store *sessions.CookieStore) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := store.Get(r, "session.id")
			if err != nil {
				ErrorJSON(w, "Internal server error", http.StatusUnauthorized)
				return
			}

			if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
				sessionId := session.Values["sessionID"]

				data, err := sess.GetToken(context.Background(), sessionId.(string))
				if err != nil {
					ErrorJSON(w, err.Error(), http.StatusForbidden)
					return
				}

				ctx := context.WithValue(r.Context(), "userID", data.UserID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				ErrorJSON(w, "Forbidden", http.StatusForbidden)
				return
			}
		})
	}
}
