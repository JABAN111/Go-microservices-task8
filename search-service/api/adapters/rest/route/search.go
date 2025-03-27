package route

import (
	"log/slog"
	"net/http"

	"yadro.com/course/api/adapters/rest/controller"
	"yadro.com/course/api/core"
)

func RegisterSearchRoutes(log *slog.Logger, mux *http.ServeMux, m core.SearchServicePort) {
	mux.HandleFunc("GET /api/search", controller.NewSearchHandler(
		log,
		m,
	))
}
