package route

import (
	"context"
	"log/slog"
	"net/http"

	"yadro.com/course/api/adapters/rest/controller"
	"yadro.com/course/api/core/ports"
)

func RegisterCommonRoutes(ctx context.Context, log *slog.Logger, mux *http.ServeMux, m ports.GrpcManager) {
	mux.HandleFunc("GET /api/ping", controller.NewPingAllHandler(
		ctx,
		log,
		m,
	))
}
