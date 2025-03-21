package grpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/core"
)

const maxAttempts = 10
const timeout = time.Second * 10

func NewServer(service core.Updater, log *slog.Logger) *Server {
	return &Server{service: service, log: log}
}

type Server struct {
	updatepb.UnimplementedUpdateServer //TODO: почитать зачем нам оно нам надо, без него ж тоже работает
	log                                *slog.Logger
	service                            core.Updater
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Status(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatusReply, error) {
	status := s.service.Status(ctx)
	var protoStatus updatepb.Status

	if status == "running" {
		protoStatus = updatepb.Status_STATUS_RUNNING
	} else {
		protoStatus = updatepb.Status_STATUS_IDLE
	}

	return &updatepb.StatusReply{Status: protoStatus}, nil
}

func (s *Server) Update(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Update(ctx); err != nil {
		s.log.Error("Failed to update", "error", err)
		return nil, err
	}
	return nil, nil
}

func (s *Server) Stats(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatsReply, error) {
	err := s.waitUpdating(ctx)
	if err != nil {
		s.log.Error("Service is still running, cannot drop")
		return nil, err
	}

	stats, err := s.service.Stats(ctx)
	if err != nil {
		s.log.Error("Failed to stats", "error", err)
		return nil, err
	}

	return &updatepb.StatsReply{
		WordsTotal:    int64(stats.WordsTotal),
		WordsUnique:   int64(stats.WordsUnique),
		ComicsTotal:   int64(stats.ComicsTotal),
		ComicsFetched: int64(stats.ComicsFetched),
	}, nil
}

func (s *Server) Drop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	err := s.waitUpdating(ctx)
	if err != nil {
		s.log.Error("Service is still running, cannot drop")
		return nil, err
	}
	if err := s.service.Drop(ctx); err != nil {
		s.log.Error("Fail to drop", "error", err)
		return nil, err
	}

	s.log.Info("Database successfully dropped")
	return &emptypb.Empty{}, nil
}

func (s *Server) waitUpdating(ctx context.Context) error {
	for i := 0; i < maxAttempts; i++ {
		status := s.service.Status(ctx)
		if status != core.StatusRunning {
			break
		}
		s.log.Info("Waiting for service to become idle before dropping...")
		time.Sleep(timeout)
	}

	if status := s.service.Status(ctx); status == core.StatusRunning {
		return errors.New("service is still running")
	}
	return nil
}
