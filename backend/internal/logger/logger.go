package logger

import (
	"log/slog"
	"os"
)

// New creates the application logger.
func New() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
}
