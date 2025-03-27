package controller

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"yadro.com/course/api/core"

	util "yadro.com/course/api/internal/utils/rest"
)

type normalizeResponse struct {
	Words []string `json:"words"`
	Total int      `json:"total"`
}

func NewPingHandler(log *slog.Logger, wordClient core.WordsServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := wordClient.Ping(r.Context()); err != nil {
			util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "error has occured while pinging word client")
			return
		}

		util.WriteResponse(r.Context(), log, w, http.StatusOK, "Pong")
	}
}

func NewNormalizeHandler(log *slog.Logger, wordClient core.WordsServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := strings.TrimSpace(r.URL.Query().Get("phrase"))

		if phrase == "" {
			util.WriteResponse(r.Context(), log, w, http.StatusBadRequest, "variable phrase is required")
			return
		}

		result, err := wordClient.Norm(r.Context(), phrase)
		if err != nil {
			if errors.Is(err, core.ErrResourceExhausted) {
				util.WriteResponseJSON(r.Context(), log, w, http.StatusBadRequest, nil)
				return
			}
			util.WriteResponseJSON(r.Context(), log, w, http.StatusInternalServerError, map[string]string{"error": "failed to normalize phrase"})
			return
		}

		response := normalizeResponse{
			Words: result,
			Total: len(result),
		}

		util.WriteResponseJSON(r.Context(), log, w, http.StatusOK, response)
	}
}
