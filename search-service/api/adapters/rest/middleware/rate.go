package middleware

import (
	"net/http"

	"yadro.com/course/api/adapters/rest"
	"yadro.com/course/api/core"
)

func Rate(limiter core.RateLimiter) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if limiter.Wait(r.Context()) != nil {
				rest.WriteError(r.Context(), w, http.StatusInternalServerError, "Internal server error, try later")
				return
			}
			next(w, r)
		}
	}
}
