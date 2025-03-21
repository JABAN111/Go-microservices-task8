package route

import (
	"context"
	"log/slog"
	"net/http"

	"yadro.com/course/api/adapters/rest/controller"
	"yadro.com/course/api/core/ports"
)

func RegisterUpdateRoutes(ctx context.Context, log *slog.Logger, mux *http.ServeMux, updateClient ports.UpdateServicePort) {
	mux.HandleFunc("GET /api/db/stats", controller.NewStatsHandler(
		ctx,
		log,
		updateClient))
	mux.HandleFunc("GET /api/db/status", controller.NewStatusHandler(
		ctx,
		log,
		updateClient,
	))

	mux.HandleFunc("POST /api/db/update", controller.NewUpdateHandler(
		ctx,
		log,
		updateClient,
	))

	mux.HandleFunc("DELETE /api/db", controller.NewDropHandler(
		ctx,
		log,
		updateClient,
	))

}
