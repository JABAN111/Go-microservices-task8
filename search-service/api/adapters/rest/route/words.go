package route

import (
	"log/slog"
	"net/http"

	"yadro.com/course/api/core"

	"yadro.com/course/api/adapters/rest/controller"
)

func RegisterNormRoutes(log *slog.Logger, mux *http.ServeMux, wordClient core.WordsServicePort) {
	mux.HandleFunc("GET /api/words/ping", controller.NewPingHandler(log, wordClient))
	mux.HandleFunc("GET /api/words", controller.NewNormalizeHandler(log, wordClient))

}
