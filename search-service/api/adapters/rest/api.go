package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core"
)

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

type errWrapper struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func WriteError(ctx context.Context, w http.ResponseWriter, httpCode int, message string) {
	WriteResponseJSON(ctx, w, httpCode, errWrapper{
		Message: message,
		Time:    time.Now(),
	})

}

var logger = slog.Default()

func WriteResponse(_ context.Context, w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain")

	var err error
	if _, err = fmt.Fprintln(w, message); err != nil {
		logger.Error("Failed to write response", "error", err)
	}

	logger.Debug("Finished sending response",
		"error", err,
		"data", message,
		"status", statusCode,
	)
}

func WriteResponseJSON(_ context.Context, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var err error
	if err = json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, writeErr := fmt.Fprintln(w, `{"error": "Internal server error"}`); writeErr != nil {
			logger.Error("Failed to write fallback error response", "error", writeErr)
		}
		logger.Error("Failed to encode JSON response", "error", err)
		return
	}

	logger.Debug("Finished sending JSON response",
		"error", err,
		"data", data,
		"status", statusCode,
	)
}

// -- Middleware

func Chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
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
		WriteResponseJSON(r.Context(), w, http.StatusOK, reply)
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
		WriteResponse(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	phraseParam, err := url.QueryUnescape(r.URL.Query().Get("phrase"))
	if err != nil {
		// It returns an error if any % is not followed by two hexadecimal digits.
		WriteError(r.Context(), w, http.StatusBadRequest, "Invalid request")
		return
	}
	if strings.TrimSpace(phraseParam) == "" {
		WriteError(r.Context(), w, http.StatusBadRequest, "phrase param is required")
		return
	}

	log.Debug("Sending phrase", "phrase", phraseParam)
	comics, err := fn(r.Context(), phraseParam, limit)
	if err != nil {
		log.Error("Search failed", "error", err)
		WriteError(r.Context(), w, http.StatusInternalServerError, "Failed to perform search")
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
		WriteError(r.Context(), w, http.StatusInternalServerError, "Failed to encode response")
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
			WriteError(r.Context(), w, http.StatusInternalServerError, "Error while getting stats")
			return
		}
		WriteResponseJSON(r.Context(), w, http.StatusOK, StatsReply{
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
			WriteError(r.Context(), w, http.StatusInternalServerError, "Error while getting drop")
		}
		WriteResponse(r.Context(), w, http.StatusOK, "Database dropped")
	}
}
func NewStatusHandler(log *slog.Logger, updateClient core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updStatus, err := updateClient.Status(r.Context())
		log.Info("Current status is specified", "status", updStatus)
		if err != nil {
			WriteError(r.Context(), w, http.StatusInternalServerError, "Error while getting status")
			return
		}
		reply := updateStatus{Status: string(updStatus)}
		if err = json.NewEncoder(w).Encode(reply); err != nil {
			log.Error("encoding failed", "error", err)
			WriteError(r.Context(), w, http.StatusInternalServerError, "Failed while writing stauts")
		}
	}
}

func NewUpdateHandler(log *slog.Logger, updateClient core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updateClient.Update(r.Context()); err != nil {
			if status.Code(err) == codes.AlreadyExists {
				WriteResponseJSON(r.Context(), w, http.StatusAccepted, "Accepted")
				return
			}
			log.Error("Internal error has been occurred", "error", err)
			WriteError(r.Context(), w, http.StatusInternalServerError, "Failed to update content")
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
			WriteError(r.Context(), w, http.StatusInternalServerError, "error has occured while pinging word client")
			return
		}

		WriteResponse(r.Context(), w, http.StatusOK, "Pong")
	}
}

func NewNormalizeHandler(log *slog.Logger, wordClient core.Worder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := strings.TrimSpace(r.URL.Query().Get("phrase"))

		if phrase == "" {
			WriteError(r.Context(), w, http.StatusBadRequest, "variable phrase is required")
			return
		}

		result, err := wordClient.Norm(r.Context(), phrase)
		if err != nil {
			if errors.Is(err, core.ErrResourceExhausted) {
				WriteError(r.Context(), w, http.StatusBadRequest, "Failed to normalize, because your request is too big")
				return
			}
			WriteError(r.Context(), w, http.StatusInternalServerError, "failed to normalize phrase")
			return
		}

		response := normalizeResponse{
			Words: result,
			Total: len(result),
		}

		WriteResponseJSON(r.Context(), w, http.StatusOK, response)
	}
}

// aaa

type User struct {
	Login    string `json:"name"`
	Password string `json:"password"`
}

func NewLoginHandler(log *slog.Logger, auth core.Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Error("failed to decode login request", slog.String("error", err.Error()))
			WriteError(r.Context(), w, http.StatusBadRequest, "bad request")
			return
		}

		token, err := auth.Login(user.Login, user.Password)
		if err != nil {
			log.Warn("unauthorized login attempt", slog.String("login", user.Login), slog.String("error", err.Error()))
			WriteError(r.Context(), w, http.StatusUnauthorized, "unauthorized")
			return
		}

		log.Info("successful login", slog.String("login", user.Login))
		WriteResponse(r.Context(), w, http.StatusOK, token)
	}
}

// Method to normalize 404 errors for common
// Example of output:
//
//	{
//	    "message": "Page not found",
//	    "time": "2025-04-06T01:51:07.887495+03:00"
//	}
func NewNotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		WriteError(r.Context(), w, http.StatusNotFound, "Page not found")
	}
}
