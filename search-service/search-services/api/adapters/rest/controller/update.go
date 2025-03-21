package controller

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"yadro.com/course/api/core/ports"
	util "yadro.com/course/api/internal/utils/rest"
)

type StatsReply struct {
	WordsTotal    int64 `json:"words_total"`
	WordsUnique   int64 `json:"words_unique"`
	ComicsTotal   int64 `json:"comics_total"`
	ComicsFetched int64 `json:"comics_fetched"`
}

const updateRunning = "STATUS_RUNNING"

func NewStatsHandler(ctx context.Context, log *slog.Logger, updateClient ports.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := updateClient.Stats(r.Context())
		if err != nil {
			if status.Code(err) == codes.FailedPrecondition {
				util.WriteResponse(ctx, log, w, http.StatusBadRequest, "Waiting service to finish updating")
				return
			}
			util.WriteResponse(ctx, log, w, http.StatusInternalServerError, "Error while getting stats")
			return
		}
		util.WriteResponseJSON(ctx, log, w, http.StatusOK, StatsReply{
			WordsTotal:    stats.WordsTotal,
			WordsUnique:   stats.WordsUnique,
			ComicsTotal:   stats.ComicsTotal,
			ComicsFetched: stats.ComicsFetched,
		})
	}
}

func NewDropHandler(ctx context.Context, log *slog.Logger, updateClient ports.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := updateClient.Drop(r.Context())
		if err != nil {
			util.WriteResponse(ctx, log, w, http.StatusInternalServerError, "Error while getting drop")
		}
		util.WriteResponse(ctx, log, w, http.StatusOK, "Database dropped") // NOTE: классная ручка
	}
}
func NewStatusHandler(ctx context.Context, log *slog.Logger, updateClient ports.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updStatus, err := updateClient.Status(context.Background())
		if err != nil {
			util.WriteResponse(ctx, log, w, http.StatusInternalServerError, "Error while getting status")
			return
		}
		util.WriteResponseJSON(ctx, log, w, http.StatusOK, getStatusMap(updStatus))
	}
}

// TODO пример ручки с ошибкой
// func NewUpdateHandler(ctx context.Context, log *slog.Logger, updateClient ports.UpdateServicePort) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		updateStatus, err := updateClient.Status(r.Context())
//		log.Debug("Got updateStatus", "updateStatus", updateStatus)
//		if err != nil {
//			log.Error("Error while getting updateStatus", "error", err)
//			util.WriteResponseJSON(ctx, log, w, http.StatusAccepted, "Error while getting update")
//			return
//		}
//
//		st := getStatusMap(updateStatus)["updateStatus"]
//		if st == "running" {
//			util.WriteResponseJSON(ctx, log, w, http.StatusAccepted, nil)
//			return
//		}
//
//		err = updateClient.Update(r.Context())
//		if err != nil {
//			log.Error("Potential error while updating updateStatus", "error", err.Error())
//			fmt.Println(err.Error())
//			st1 := codes.FailedPrecondition
//			st2 := status.Code(err)
//			log.Info("Got st", "st1", st1, "st2", st2, "st1==st2", st1 == st2)
//			if status.Code(err) == codes.FailedPrecondition {
//
//				util.WriteResponseJSON(ctx, log, w, http.StatusAccepted, "")
//				return
//			}
//			util.WriteResponseJSON(ctx, log, w, http.StatusAccepted, "Error while getting update")
//			return
//		}
//		log.Debug("Reply ok", "status", st)
//		w.WriteHeader(http.StatusOK)
//		//util.WriteResponseJSON(ctx, log, w, http.StatusOK, "")
//	}
//}

func NewUpdateHandler(_ context.Context, log *slog.Logger, updateClient ports.UpdateServicePort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updStatus, err := getStringStatus(updateClient, log)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// TODO: привязаться бы к константам из другого микросервиса
		if updStatus == updateRunning {
			log.Debug("Service is already running")
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// WTF: без горутины, то есть в синхронном режиме падали тесты после отдачи респонса.
		go func() {
			log.Debug("Started updating")
			if err = updateClient.Update(context.Background()); err != nil {
				log.Error("Error while updating", "error", err)
				return
			}
			log.Debug("Finished updating")
		}()

		log.Debug("Updating finished", "error", err)
		if err != nil {
			if status.Code(err) == codes.FailedPrecondition {
				log.Debug("Already updating")
				w.WriteHeader(http.StatusAccepted)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug("Reply ok", "status", updStatus)
		w.WriteHeader(http.StatusOK)               // Сигнализирует о начале обновления, но не о его "удачном" завершении
		log.Debug("Reply ok", "status", updStatus) //NOTE: по логам это печатается
	}
}

func getStatusMap(status string) map[string]string {
	stat := strings.Split(status, ":")[1]
	return map[string]string{"status": strings.ToLower(strings.Split(stat, "_")[1])}
}

func getStringStatus(updateClient ports.UpdateServicePort, log *slog.Logger) (string, error) {
	updateStatus, err := updateClient.Status(context.Background())

	if err != nil {
		log.Error("Error while getting updateStatus", "error", err)
		return "", err
	}

	updStatus := getStatusMap(updateStatus)["status"]
	return updStatus, nil
}
