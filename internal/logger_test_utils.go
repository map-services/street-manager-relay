package internal

import (
	"bytes"
	"log/slog"
)

// SetupSlogBuffer returns a buffer and a logger that writes JSON to it.
func SetupSlogBuffer() (*bytes.Buffer, *slog.Logger) {
	buf := new(bytes.Buffer)
	handler := slog.NewJSONHandler(buf, nil)
	return buf, slog.New(handler)
}
