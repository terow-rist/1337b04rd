package httpserver

import (
	"1337bo4rd/internal/core/domain"
	"1337bo4rd/internal/core/service"
	"context"
	"log/slog"
	"net/http"
)

type ctxKey string

const sessionKey ctxKey = "session"

func SessionMiddleware(user *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Info("http request", "method", r.Method, "path", r.URL.Path)

			var sid string
			if c, err := r.Cookie("session_id"); err == nil {
				sid = c.Value
			}
			sess, isNew, err := user.FindOrCreate(sid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if isNew {
				slog.Info("new session", "id", sess.ID)
				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    sess.ID,
					Path:     "/",
					Expires:  sess.ExpiresAt,
					HttpOnly: true,
					SameSite: http.SameSiteLaxMode,
					Secure:   true,
				})
			}
			ctx := context.WithValue(r.Context(), sessionKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getSession(r *http.Request) *domain.User {
	if v := r.Context().Value(sessionKey); v != nil {
		return v.(*domain.User)
	}
	return nil
}
