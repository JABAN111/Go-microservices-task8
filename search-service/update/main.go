package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"

	updatepb "yadro.com/course/proto/update"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"yadro.com/course/update/adapters/db"
	updategrpc "yadro.com/course/update/adapters/grpc"
	"yadro.com/course/update/adapters/words"
	"yadro.com/course/update/adapters/xkcd"
	"yadro.com/course/update/config"
	"yadro.com/course/update/core"
)

const defaultPath = "config.yaml"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var configPath string
	flag.StringVar(&configPath, "config", defaultPath, "server configuration file")
	flag.Parse()
	cfg := config.MustLoad(configPath)

	log := mustMakeLogger(cfg.LogLevel)

	log.Info("starting server")
	log.Debug("debug messages are enabled")

	log.Info("Concurrency setup", "concurrency", cfg.XKCD.Concurrency)
	storage, err := db.New(log, cfg.DBAddress, cfg.BatchSize, cfg.XKCD.Concurrency)
	if err != nil {
		log.Error("failed to connect to db", "error", err)
		stop()
		return
	}
	if err := storage.Migrate(); err != nil {
		log.Error("failed to migrate db", "error", err)
		stop()
		return
	}

	xkcdClient, err := xkcd.NewClient(cfg.XKCD.URL, cfg.XKCD.Timeout, log)
	if err != nil {
		log.Error("failed create XKCD client", "error", err)
		stop()
		return
	}

	wordsClient, err := words.NewClient(cfg.WordsAddress, log)
	if err != nil {
		log.Error("failed create Words client", "error", err)
		stop()
		return
	}

	updater, err := core.NewService(log, storage, xkcdClient, wordsClient, cfg.XKCD.Concurrency)
	if err != nil {
		log.Error("failed create Update service", "error", err)
		stop()
		return
	}

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error("failed to listen", "error", err)
		stop()
		return
	}
	log.Info("server listening on", "address", cfg.Address)

	s := grpc.NewServer()
	updatepb.RegisterUpdateServer(s, updategrpc.NewServer(updater, log))
	reflection.Register(s)

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		s.GracefulStop()
	}()

	if err := s.Serve(listener); err != nil {
		log.Error("failed to serve", "error", err)
		stop()
		return
	}
}

func mustMakeLogger(logLevel string) *slog.Logger {
	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "ERROR":
		level = slog.LevelError
	default:
		panic("unknown log level: " + logLevel)
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
