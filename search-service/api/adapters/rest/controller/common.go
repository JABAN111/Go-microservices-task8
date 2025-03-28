package controller

import (
	"log/slog"
	"net/http"

	"yadro.com/course/api/core"

	util "yadro.com/course/api/internal/utils/rest"
)

type pingResponse struct {
	Replies map[string]string `json:"replies"`
}

func NewPingAllHandler(log *slog.Logger, m core.GrpcManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceMap := m.PingAll(r.Context())
		response := pingResponse{Replies: serviceMap}

		log.Debug("Client ping services", "response", response)
		util.WriteResponseJSON(r.Context(), log, w, http.StatusOK, response)
	}
}
