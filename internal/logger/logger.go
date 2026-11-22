package logger

import (
	"log/slog"
	"os"
)

//This function is a custom/self-made logger, this logger returns json
func New(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(handler)
}