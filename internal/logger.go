package internal

import (
	"log/slog"
	"os"
)

// SetupLogger configures the default structured JSON logger.
func SetupLogger() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
