package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
	util "yadro.com/course/api/internal/utils/rest"
)

type StatsReply struct {
	WordsTotal    int64 `json:"words_total"`
	WordsUnique   int64 `json:"words_unique"`
	ComicsTotal   int64 `json:"comics_total"`
	ComicsFetched int64 `json:"comics_fetched"`
}
type updateStatus struct {
	Status string `json:"status"`
}

func NewStatsHandler(log *slog.Logger, updateClient core.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := updateClient.Stats(r.Context())
		if err != nil {
			util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "Error while getting stats")
			return
		}
		util.WriteResponseJSON(r.Context(), log, w, http.StatusOK, StatsReply{
			WordsTotal:    stats.WordsTotal,
			WordsUnique:   stats.WordsUnique,
			ComicsTotal:   stats.ComicsTotal,
			ComicsFetched: stats.ComicsFetched,
		})
	}
}

func NewDropHandler(log *slog.Logger, updateClient core.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updateClient.Drop(r.Context())
		if err != nil {
			util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "Error while getting drop")
		}
		util.WriteResponse(r.Context(), log, w, http.StatusOK, "Database dropped") // NOTE: классная ручка
	}
}
func NewStatusHandler(log *slog.Logger, updateClient core.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updStatus, err := updateClient.Status(r.Context())
		if err != nil {
			util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "Error while getting status")
			return
		}
		reply := updateStatus{Status: string(updStatus)}
		if err = json.NewEncoder(w).Encode(reply); err != nil {
			log.Error("encoding failed", "error", err)
		}
		util.WriteResponseJSON(r.Context(), log, w, http.StatusOK, reply)
	}
}

func NewUpdateHandler(log *slog.Logger, updateClient core.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updateClient.Update(r.Context()); err != nil {
			if status.Code(err) == codes.AlreadyExists {
				w.WriteHeader(http.StatusAccepted)
				return
			}
			log.Error("Internal error has been occurred", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
