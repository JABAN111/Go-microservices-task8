package main

import (
	"context"
	"flag"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	searchpb "yadro.com/course/proto/search"
	"yadro.com/course/search/adapters/db"

	searchgrpc "yadro.com/course/search/adapters/grpc"
	wordsgrpc "yadro.com/course/search/adapters/words"
	"yadro.com/course/search/config"
	"yadro.com/course/search/core"
	"yadro.com/course/search/logger"
)

const defaultPath = "config.yaml"

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", defaultPath, "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.LogLevel)
	log.Debug("Debug level is enabled", "cfg", cfg)
	ctx := context.Background()

	wordClient, err := wordsgrpc.NewClient(cfg.WordsAddress, log)

	if err != nil {
		log.Error("Fail to init connection to word microservice")
		panic(err)
	}
	if err = wordClient.Ping(context.Background()); err != nil {
		panic(err)
	}
	dbClient, err := db.New(log, cfg.DBAddress, cfg.Workers, cfg.DBMaxConnections, cfg.DBMaxLifeTime)

	if err != nil {
		log.Error("Fail to init connection to database", "db address", cfg.DBAddress)
		panic(err)
	}

	service := core.NewService(ctx, log, dbClient, cfg.IndexTtl, wordClient, cfg.Workers)

	// grpc server
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}
	log.Info("server listening on", "address", cfg.Address)

	s := grpc.NewServer()
	searchpb.RegisterSearchServer(s, searchgrpc.NewServer(log, service))
	reflection.Register(s)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Debug("shutting down server")
		s.GracefulStop()
	}()

	if err := s.Serve(listener); err != nil {
		log.Error("failed to serve", "error", err)
	}

}
