package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"yadro.com/course/api/adapters/grpc/client"
	"yadro.com/course/api/adapters/grpc/managers"
	"yadro.com/course/api/adapters/rest"
	"yadro.com/course/api/adapters/rest/route"
	"yadro.com/course/api/internal/logger"

	"yadro.com/course/api/internal/config"
)

const maxShutdownTime = 8 * time.Second

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()
	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.LogLevel)
	log.Debug("Debug level is enabled")
	log.Info("Config parsed", "config", cfg)

	clientManager := managers.NewClientManager(log)
	wordsClient, err := client.NewWordsClient(cfg.WordsAddress, log)
	if err != nil {
		log.Error("Cannot init WordsClient", "error", err)
		os.Exit(1)
	}
	updateClient, err := client.NewUpdateClient(cfg.UpdateAddress, log)
	if err != nil {
		log.Error("Cannot init UpdateClient", "error", err)
		os.Exit(1)
	}
	searchClient, err := client.NewSearchClient(cfg.SearchAddress, log)
	if err != nil {
		log.Error("Cannot init search client", "error", err)
		os.Exit(1)
	}

	clientManager.Register("words", wordsClient)
	log.Info("Words client successfully registered")
	clientManager.Register("update", updateClient)
	log.Info("Update client successfully registered")
	clientManager.Register("search", searchClient)
	log.Debug("Status of clients", "status", clientManager.PingAll(context.Background()))

	rootMux := http.NewServeMux()

	route.RegisterCommonRoutes(log, rootMux, clientManager)
	route.RegisterNormRoutes(log, rootMux, wordsClient)
	route.RegisterUpdateRoutes(log, rootMux, updateClient)
	route.RegisterSearchRoutes(log, rootMux, searchClient)

	restServer := rest.NewServer(log, rootMux, cfg.HttpServer.ServerAddress, cfg.HttpServer.HttpTimeout)
	log.Info("Server started", "config", cfg)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := restServer.Run(); err != nil {
			log.Error("Server failed", "error", err)
			gracefulShutdown(log, restServer, clientManager)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	gracefulShutdown(log, restServer, clientManager)
}

func gracefulShutdown(log *slog.Logger, restServer *rest.Server, clientManager *managers.ClientManager) {
	log.Info("Shutting down the server")
	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), maxShutdownTime)
	defer cancel()

	log.Debug("Closing all clients...")
	clientManager.CloseAll(shutdownTimeoutCtx)
	log.Debug("Client closing are finished")

	log.Debug("Starting shutdown for the http server")
	if err := restServer.Http.Shutdown(shutdownTimeoutCtx); err != nil {
		log.Error("Failed to shut down server", "error", err)
	}

	log.Info("Server shutdown complete")
}
