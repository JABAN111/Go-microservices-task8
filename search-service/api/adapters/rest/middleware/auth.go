package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"yadro.com/course/api/adapters/rest"
	"yadro.com/course/api/core"
)

func Auth(verifier core.TokenVerifier) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			fmt.Printf("header: %s", authHeader)
			authFields := strings.Fields(authHeader)

			if len(authFields) != 2 || strings.ToLower(authFields[0]) != "bearer" {
				rest.WriteError(r.Context(), w, http.StatusUnauthorized, "Invalid token")
				return
			}

			token := authFields[1]

			if err := verifier.Verify(token); err != nil {
				rest.WriteError(r.Context(), w, http.StatusUnauthorized, "Invalid token")
				return
			}

			next(w, r)
		}
	}
}
