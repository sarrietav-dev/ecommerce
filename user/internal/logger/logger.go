package logger

import (
	"log/slog"
	"net"
	"os"
)

var Logger *slog.Logger

func InitLogger() {
	conn, err := net.Dial("tcp", "localhost:5044")
	if err != nil {
		Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, // Set log level
		}))
		Logger.Warn("Failed to connect to Logstash", slog.String("error", err.Error()))
		return
	}

	Logger = slog.New(slog.NewJSONHandler(conn, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
