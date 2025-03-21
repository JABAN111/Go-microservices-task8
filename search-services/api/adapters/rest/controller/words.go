package controller

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	ownErr "yadro.com/course/api/core/errors"
	"yadro.com/course/api/core/ports"
	util "yadro.com/course/api/internal/utils/rest"
)

type NormalizeResponse struct {
	Words []string `json:"words"`
	Total int      `json:"total"`
}

func NewPingHandler(ctx context.Context, log *slog.Logger, wordClient ports.WordsServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := wordClient.Ping(ctx); err != nil {
			util.WriteResponse(ctx, log, w, http.StatusInternalServerError, "error has occured while pinging word client")
			return
		}

		util.WriteResponse(ctx, log, w, http.StatusOK, "Pong")
	}
}

func NewNormalizeHandler(ctx context.Context, log *slog.Logger, wordClient ports.WordsServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := strings.TrimSpace(r.URL.Query().Get("phrase"))

		if phrase == "" {
			util.WriteResponse(ctx, log, w, http.StatusBadRequest, "variable phrase is required")
			return
		}

		result, err := wordClient.Norm(ctx, phrase)
		if err != nil {
			if errors.Is(err, ownErr.ErrResourceExhausted) {
				util.WriteResponseJSON(ctx, log, w, http.StatusBadRequest, nil)
				return
			}
			util.WriteResponseJSON(ctx, log, w, http.StatusInternalServerError, map[string]string{"error": "failed to normalize phrase"})
			return
		}

		response := NormalizeResponse{
			Words: result,
			Total: len(result),
		}

		util.WriteResponseJSON(ctx, log, w, http.StatusOK, response)
	}
}
