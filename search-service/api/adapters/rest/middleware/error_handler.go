package middleware

import (
	"context"
	"log"
	"net/http"

	"yadro.com/course/api/adapters/rest"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v", r)
				rest.WriteError(context.Background(), w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
