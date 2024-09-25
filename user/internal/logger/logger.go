package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
