package logger

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var (
	instance *slog.Logger
	once     sync.Once
)

const (
	//NOTE: в чистом виде баловство)
	envDev  = "dev"
	envProd = "prod"
)

func GetInstance() *slog.Logger {
	once.Do(func() {
		env := os.Getenv("LOG_ENV")
		switch env {
		case envDev:
			instance = setupLogger(envDev)
		case envProd:
			instance = setupLogger(envProd)
		default:
			fmt.Println("Error: Failed to initialize logger, using default dev logger.")
			instance = setupLogger(envDev)
		}
	})
	return instance
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {

	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
