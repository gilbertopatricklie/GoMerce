package logger

import (
	"log/slog"
	"os"
)

//This function is a custom/self-made logger, this logger returns json
func New() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})

	return slog.New(handler)
}