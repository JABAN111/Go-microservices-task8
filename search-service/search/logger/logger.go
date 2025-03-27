package logger

import (
	"log/slog"
	"os"
)

func MustMakeLogger(logLevel string) *slog.Logger {
	levelMap := map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}

	level, ok := levelMap[logLevel]
	if !ok {
		return slog.Default()
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
