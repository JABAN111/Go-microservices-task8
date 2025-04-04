package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
	"yadro.com/course/api/util"
)

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}
type Authenticator interface {
	Login(user, password string) (string, error)
}

func NewPingAllHandler(log *slog.Logger, clients map[string]core.GrpcClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := PingResponse{
			Replies: make(map[string]string, len(clients)),
		}
		for name, client := range clients {
			if err := client.Ping(r.Context()); err != nil {
				reply.Replies[name] = "unavailable"
				log.Error("one of services is not available", "service", name, "error", err)
				continue
			}
			reply.Replies[name] = "ok"
		}
		util.WriteResponseJSON(r.Context(), log, w, http.StatusOK, reply)
	}
}

// --- Searcher methods
type comicsReply struct {
	Comics []comic `json:"comics"`
	Total  int     `json:"total"`
}

type comic struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

type searchFunc func(ctx context.Context, phrase string, limit int64) ([]core.Comics, error)

func handleSearch(log *slog.Logger, fn searchFunc, w http.ResponseWriter, r *http.Request) {
	lp := strings.TrimSpace(r.URL.Query().Get("limit"))
	limit := int64(10)

	limit, err := parseLimit(lp, limit)
	if err != nil {
		util.WriteResponse(r.Context(), log, w, http.StatusBadRequest, err.Error())
		return
	}

	phraseParam, err := url.QueryUnescape(r.URL.Query().Get("phrase"))
	if err != nil {
		// It returns an error if any % is not followed by two hexadecimal digits.
		util.WriteResponse(r.Context(), log, w, http.StatusBadRequest, "Invalid request")
		return
	}
	if strings.TrimSpace(phraseParam) == "" {
		util.WriteResponse(r.Context(), log, w, http.StatusBadRequest, "phrase param is required")
		return
	}

	log.Debug("Sending phrase", "phrase", phraseParam)
	comics, err := fn(r.Context(), phraseParam, limit)
	if err != nil {
		log.Error("Search failed", "error", err)
		util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "Failed to perform search")
		return
	}

	comicRep := make([]comic, len(comics))
	for i, com := range comics {
		id := com.ID
		comicRep[i] = comic{
			ID:  id,
			Url: com.ImgUrl,
		}
	}

	response := comicsReply{
		Comics: comicRep,
		Total:  len(comicRep),
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Error("encoding failed", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "Failed to encode response")
		return
	}

	log.Debug("Finished searching", "result", comics)

}

func NewSearchHandler(log *slog.Logger, client core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleSearch(log, client.Search, w, r)
	}
}
func NewISearchHandler(log *slog.Logger, client core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleSearch(log, client.ISearch, w, r)
	}
}

func parseLimit(param string, defaultValue int64) (int64, error) {
	if param == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		if errors.As(err, new(*strconv.NumError)) {
			return 0, core.ErrInvalidType
		}
		return 0, core.ErrFailedToProcessLimit
	}

	return parsed, nil
}

// --- Updater methods
type StatsReply struct {
	WordsTotal    int64 `json:"words_total"`
	WordsUnique   int64 `json:"words_unique"`
	ComicsTotal   int64 `json:"comics_total"`
	ComicsFetched int64 `json:"comics_fetched"`
}
type updateStatus struct {
	Status string `json:"status"`
}

func NewStatsHandler(log *slog.Logger, updateClient core.Updater) http.HandlerFunc {
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

func NewDropHandler(log *slog.Logger, updateClient core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updateClient.Drop(r.Context())
		if err != nil {
			util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "Error while getting drop")
		}
		util.WriteResponse(r.Context(), log, w, http.StatusOK, "Database dropped") // NOTE: классная ручка
	}
}
func NewStatusHandler(log *slog.Logger, updateClient core.Updater) http.HandlerFunc {
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

func NewUpdateHandler(log *slog.Logger, updateClient core.Updater) http.HandlerFunc {
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

// --- words

type normalizeResponse struct {
	Words []string `json:"words"`
	Total int      `json:"total"`
}

func NewPingHandler(log *slog.Logger, wordClient core.Worder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := wordClient.Ping(r.Context()); err != nil {
			util.WriteResponse(r.Context(), log, w, http.StatusInternalServerError, "error has occured while pinging word client")
			return
		}

		util.WriteResponse(r.Context(), log, w, http.StatusOK, "Pong")
	}
}

func NewNormalizeHandler(log *slog.Logger, wordClient core.Worder) http.HandlerFunc {
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
