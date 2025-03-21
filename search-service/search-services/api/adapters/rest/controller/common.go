package controller

import (
	"context"
	"log/slog"
	"net/http"

	"yadro.com/course/api/core/ports"
	util "yadro.com/course/api/internal/utils/rest"
)

type pingResponse struct {
	Replies map[string]string `json:"replies"`
}

func NewPingAllHandler(ctx context.Context, log *slog.Logger, m ports.GrpcManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceMap := m.PingAll(ctx)
		response := pingResponse{Replies: serviceMap}

		log.Debug("Client ping services", "response", response)
		util.WriteResponseJSON(ctx, log, w, http.StatusOK, response)
	}
}
