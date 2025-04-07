package main

import (
	"context"
	"flag"
	"net/http"
	"sync"
	"time"

	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"yadro.com/course/api/adapters/aaa"
	grpc "yadro.com/course/api/adapters/grpc"
	"yadro.com/course/api/adapters/rest"
	"yadro.com/course/api/adapters/rest/middleware"
	"yadro.com/course/api/config"
	"yadro.com/course/api/core"
	"yadro.com/course/api/logger"
)

const maxShutdownTime = 8 * time.Second

// grpc client name/tags
const (
	search = "search"
	words  = "words"
	update = "update"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()
	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.LogLevel)
	log.Debug("Debug level is enabled")
	log.Info("Config parsed", "config", cfg)

	wordsClient, err := grpc.NewWordsClient(cfg.WordsAddress, log)
	if err != nil {
		log.Error("Cannot init WordsClient", "error", err)
		os.Exit(1)
	}
	updateClient, err := grpc.NewUpdateClient(cfg.UpdateAddress, log)
	if err != nil {
		log.Error("Cannot init UpdateClient", "error", err)
		os.Exit(1)
	}
	searchClient, err := grpc.NewSearchClient(cfg.SearchAddress, log)
	if err != nil {
		log.Error("Cannot init search client", "error", err)
		os.Exit(1)
	}

	clients := map[string]core.GrpcClient{words: wordsClient, update: updateClient, search: searchClient}

	mux := http.NewServeMux()
	mux.HandleFunc("/", rest.NewNotFoundHandler())

	if err != nil {
		log.Error("Failed to init auth provider")
	}

	ctx := context.Background()
	authProvider, err := aaa.New(cfg.TokenTTL, log)
	if err != nil {
		log.Error("Failed to init auth provider", "config", cfg)
		os.Exit(1)
	}
	// -- login
	mux.HandleFunc("POST /api/login", rest.Chain(rest.NewLoginHandler(log, authProvider), middleware.Loger(log)))

	// -- common
	mux.HandleFunc("GET /api/ping", rest.NewPingAllHandler(log, clients))

	// -- search
	mux.HandleFunc("GET /api/search", rest.Chain(
		rest.NewSearchHandler(log, searchClient),
		middleware.Concurrency(cfg.ConcurrencyLimiter),
	))
	mux.HandleFunc("GET /api/isearch", rest.Chain(
		rest.NewISearchHandler(log, searchClient),
		middleware.Rate(rest.NewRateLimiter(ctx, cfg.RateLimiter))))

	// -- update
	mux.HandleFunc("GET /api/db/stats", rest.NewStatsHandler(log, updateClient))
	mux.HandleFunc("GET /api/db/status", rest.NewStatusHandler(log, updateClient))
	mux.HandleFunc("POST /api/db/update", rest.Chain(rest.NewUpdateHandler(log, updateClient),
		middleware.Loger(log),
		middleware.Auth(authProvider),
	))

	mux.HandleFunc("DELETE /api/db", rest.Chain(
		rest.NewDropHandler(log, updateClient),
		middleware.Loger(log),
		middleware.Auth(authProvider)))

	// -- words
	mux.HandleFunc("GET /api/words/ping", rest.NewPingHandler(log, wordsClient))
	mux.HandleFunc("GET /api/words", rest.NewNormalizeHandler(log, wordsClient))

	log.Info("Server started", "config", cfg)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:         cfg.HttpServer.ServerAddress,
		ReadTimeout:  cfg.HttpServer.HttpTimeout,
		WriteTimeout: cfg.HttpServer.HttpTimeout,
		Handler:      mux,
	}
	// один раз обжегся об это, пусть теперь будет warn
	if cfg.HttpServer.HttpTimeout < time.Second*5 {
		log.Warn("timeout is less then 5 seconds")
	}

	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("server closed unexpectedly", "error", err)
		}
	}()

	<-ctx.Done()
	gracefulShutdown(log, s, clients)
}

func gracefulShutdown(log *slog.Logger, restServer *http.Server, clients map[string]core.GrpcClient) {
	log.Info("Shutting down the server")
	shutdownTimeoutCtx, cancel := context.WithTimeout(context.Background(), maxShutdownTime)
	defer cancel()

	log.Debug("Closing all clients...")

	var wg sync.WaitGroup

	for _, client := range clients {
		wg.Add(1)
		go func(client core.GrpcClient) {
			defer wg.Done()

			done := make(chan error, 1)

			go func() {
				done <- client.Close()
			}()

			select {
			case <-shutdownTimeoutCtx.Done():
				log.Warn("Time-out of client disconnecting")
				return
			case err := <-done:
				if err != nil {
					log.Error("Error while closing the client", "error", err)
				}
				return
			}
		}(client)
	}
	wg.Wait()

	log.Debug("Client closing are finished")

	log.Debug("Starting shutdown for the http server")
	if err := restServer.Shutdown(shutdownTimeoutCtx); err != nil {
		log.Error("Failed to shut down server", "error", err)
	}

	log.Info("Server shutdown complete")
}
