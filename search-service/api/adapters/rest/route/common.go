package route

import (
	"log/slog"
	"net/http"

	"yadro.com/course/api/core"

	"yadro.com/course/api/adapters/rest/controller"
)

func RegisterCommonRoutes(log *slog.Logger, mux *http.ServeMux, m core.GrpcManager) {
	mux.HandleFunc("GET /api/ping", controller.NewPingAllHandler(log, m))
}
