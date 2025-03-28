package route

import (
	"log/slog"
	"net/http"

	"yadro.com/course/api/core"

	"yadro.com/course/api/adapters/rest/controller"
)

func RegisterUpdateRoutes(log *slog.Logger, mux *http.ServeMux, updateClient core.UpdateServicePort) {
	mux.HandleFunc("GET /api/db/stats", controller.NewStatsHandler(log, updateClient))
	mux.HandleFunc("GET /api/db/status", controller.NewStatusHandler(log, updateClient))

	mux.HandleFunc("POST /api/db/update", controller.NewUpdateHandler(log, updateClient))

	mux.HandleFunc("DELETE /api/db", controller.NewDropHandler(log, updateClient))
}
