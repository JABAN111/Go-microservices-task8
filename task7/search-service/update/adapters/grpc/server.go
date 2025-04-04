package grpc

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"
	"yadro.com/course/update/core"
)

func NewServer(service core.Updater, log *slog.Logger) *Server {
	return &Server{service: service, log: log}
}

type Server struct {
	updatepb.UnimplementedUpdateServer
	log     *slog.Logger
	service core.Updater
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *Server) Status(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatusReply, error) {
	st := s.service.Status(ctx)
	var protoStatus updatepb.Status

	if st == core.StatusRunning {
		protoStatus = updatepb.Status_STATUS_RUNNING
	} else {
		protoStatus = updatepb.Status_STATUS_IDLE
	}

	return &updatepb.StatusReply{Status: protoStatus}, nil
}

func (s *Server) Update(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Update(ctx); err != nil {
		if errors.Is(err, core.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "Already updating")
		}
		s.log.Error("Failed to update", "error", err)
		return nil, status.Error(codes.Internal, "Failed to update")
	}
	return nil, nil
}

func (s *Server) Stats(ctx context.Context, _ *emptypb.Empty) (*updatepb.StatsReply, error) {
	stats, err := s.service.Stats(ctx)
	if err != nil {
		s.log.Error("Failed to stats", "error", err)
		return nil, status.Error(codes.Internal, "Failed to getting stats")
	}

	return &updatepb.StatsReply{
		WordsTotal:    int64(stats.WordsTotal),
		WordsUnique:   int64(stats.WordsUnique),
		ComicsTotal:   int64(stats.ComicsTotal),
		ComicsFetched: int64(stats.ComicsFetched),
	}, nil
}

func (s *Server) Drop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Drop(ctx); err != nil {
		s.log.Error("Fail to drop", "error", err)
		return nil, status.Error(codes.Internal, "Fail to drop")
	}

	s.log.Info("Database successfully dropped")
	return &emptypb.Empty{}, nil
}
