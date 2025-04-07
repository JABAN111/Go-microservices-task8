package middleware

import (
	"log/slog"
	"net/http"

	"yadro.com/course/api/adapters/rest"
)

func Loger(log *slog.Logger) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Info("Request processing", "uri", r.RequestURI, "method", r.Method)
			next(w, r)
		}
	}
}
