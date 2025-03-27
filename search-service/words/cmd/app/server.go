package main

import (
	"flag"
	"net"
	"os"

	"google.golang.org/grpc/reflection"

	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"

	"yadro.com/course/words/internal/config"
	controller "yadro.com/course/words/internal/controller"
	"yadro.com/course/words/pkg/logger"
	"yadro.com/course/words/pkg/shutdown"

	wordspb "yadro.com/course/proto/words"
)

var log = logger.GetInstance()

const defaultConfigPath = "config.yaml"

var configPath string

func getConfig() *config.Config {
	cfg := config.NewConfig()

	if err := cleanenv.ReadConfig(configPath, cfg); err == nil {
		if cfg.BindAddress != "" {
			return cfg
		}
	}

	log.Warn("Failed to load configuration from file or environment variables")
	return config.DefaultConfig()
}

func main() {
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to config file")
	flag.Parse()

	cfg := getConfig()
	log.Info("Using configuration", "config", cfg)

	address := cfg.BindAddress

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Error("Failed to initilize listener", "err", err)
		os.Exit(-1)
	}

	grpcSrv := grpc.NewServer()
	wordspb.RegisterWordsServer(grpcSrv, controller.NewNormalizer())
	reflection.Register(grpcSrv)

	err = shutdown.RunWithShutdown(grpcSrv, func() error {
		log.Info("Starting gRPC server...")
		return grpcSrv.Serve(listener)
	})

	if err != nil {
		log.Error("Server error", "err", err)
		os.Exit(-1)
	}

}
