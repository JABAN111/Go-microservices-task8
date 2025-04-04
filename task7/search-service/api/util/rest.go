package util

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func WriteResponse(_ context.Context, log *slog.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	var err error
	if _, err = fmt.Fprintln(w, message); err != nil {
		log.Error("Failed to write response", "error", err)
	}
	log.Debug("Finish sending answer", "err", err, "data", message, "code", statusCode)

}

func WriteResponseJSON(_ context.Context, log *slog.Logger, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	var err error
	if err = json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := fmt.Fprintln(w, `{"error": "Internal server error"}`); err != nil {
			log.Error("Failed to write error response", "error", err)
		}
	}
	log.Debug("Finish sending answer", "err", err, "data", data, "code", statusCode)
}
