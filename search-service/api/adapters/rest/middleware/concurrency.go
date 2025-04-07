package middleware

import (
	"net/http"

	"yadro.com/course/api/adapters/rest"
)

func Concurrency(limit int) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		maxRequest := make(chan struct{}, limit)

		return func(w http.ResponseWriter, r *http.Request) {
			select {
			case maxRequest <- struct{}{}:
				defer func() { <-maxRequest }()
				next(w, r)
			default:
				rest.WriteError(r.Context(), w, http.StatusServiceUnavailable, "Service Unavailable (too many requests)")
			}
		}
	}
}
