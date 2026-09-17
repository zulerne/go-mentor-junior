package logger

import (
	"log/slog"
	"os"

	"github.com/zulerne/go-mentor-junior/order/internal/config"
)

// TODO (review): Do I need this small package? Or Its better to init it in main.go?

func New(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case config.EnvProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true}))
	}

	return log
}
