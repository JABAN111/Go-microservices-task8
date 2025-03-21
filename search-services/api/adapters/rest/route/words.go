package route

import (
	"context"
	"log/slog"
	"net/http"

	"yadro.com/course/api/adapters/rest/controller"
	"yadro.com/course/api/core/ports"
)

func RegisterNormRoutes(ctx context.Context, log *slog.Logger, mux *http.ServeMux, wordClient ports.WordsServicePort) {
	mux.HandleFunc("GET /api/words/ping", controller.NewPingHandler(
		ctx,
		log,
		wordClient,
	))
	mux.HandleFunc("GET /api/words", controller.NewNormalizeHandler(
		ctx,
		log,
		wordClient))

}
