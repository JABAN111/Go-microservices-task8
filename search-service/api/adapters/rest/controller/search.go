package controller

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"yadro.com/course/api/core"
)

type comicsReply struct {
	Comics []comic `json:"comics"`
	Total  int     `json:"total"`
}

type comic struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

func NewSearchHandler(log *slog.Logger, client core.SearchServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Got request on search")

		lp := strings.TrimSpace(r.URL.Query().Get("limit"))
		limit := int64(10)
		// TODO: избавиться от ада из if
		if lp != "" {
			parsed, err := strconv.ParseInt(lp, 10, 64)
			if err != nil {
				var numErr *strconv.NumError
				if errors.As(err, &numErr) {
					http.Error(w, "limit must be a number", http.StatusBadRequest)
				} else {
					http.Error(w, "Failed to process limit", http.StatusInternalServerError)
				}
				return
			}
			limit = parsed
		} else {
			log.Warn("User didn't specify limit, using default", "limit", limit)
		}

		phraseParam, err := url.QueryUnescape(r.URL.Query().Get("phrase"))
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if strings.TrimSpace(phraseParam) == "" {
			http.Error(w, "phrase param is required", http.StatusBadRequest)
			return
		}

		log.Debug("Sending phrase", "phrase", phraseParam)
		search, err := client.Search(r.Context(), phraseParam, limit)
		if err != nil {
			log.Error("Search failed", "error", err)
			http.Error(w, "Failed to perform search", http.StatusInternalServerError)
			return
		}

		comicRep := make([]comic, len(search.Comics))
		for i, com := range search.Comics {
			id, err := strconv.Atoi(com.ID)
			if err != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}
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
			return
		}

		log.Debug("Finished searching", "result", search)
	}
}
